package app

import (
	"log"
	"net/http"
	"strings"
	"time"

	"ment-agents/backend/internal/compiler"
	"ment-agents/backend/internal/model"
	"ment-agents/backend/internal/store"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store    *store.Store
	compiler *compiler.Service
}

func NewHandler(s *store.Store, c *compiler.Service) *Handler {
	return &Handler{store: s, compiler: c}
}

type createProjectRequest struct {
	Name string `json:"name"`
}

type chatRequest struct {
	Message string `json:"message"`
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.Use(cors())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		api.GET("/projects", h.listProjects)
		api.POST("/projects", h.createProject)
		api.GET("/projects/:id", h.getProject)

		api.POST("/projects/:id/chat", h.chat)
		api.GET("/projects/:id/chat", h.getChat)
		api.POST("/projects/:id/compile", h.compile)
		api.POST("/projects/:id/assemble", h.assemble)
		api.GET("/projects/:id/artifacts", h.artifacts)
	}
}

func (h *Handler) listProjects(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"projects": h.store.ListProjects()})
}

func (h *Handler) createProject(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目名称不能为空"})
		return
	}

	project := h.store.CreateProject(name)
	c.JSON(http.StatusCreated, gin.H{"project": project})
}

func (h *Handler) getProject(c *gin.Context) {
	project, err := h.store.GetProject(c.Param("id"))
	if err != nil {
		h.notFound(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *Handler) chat(c *gin.Context) {
	projectID := c.Param("id")

	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}

	history, err := h.store.GetConversation(projectID)
	if err != nil {
		h.notFound(c)
		return
	}

	log.Printf("[chat] project=%s user_message=%q", projectID, message)

	clarified, clarifyErr := h.compiler.Clarify(toClarifyInput(projectID, history, message))
	if clarifyErr != nil {
		log.Printf("[chat] clarification error=%v", clarifyErr)
	}

	agentReply := strings.TrimSpace(clarified.UserReply)
	if agentReply == "" {
		agentReply = h.buildAgentReply(message)
	}

	log.Printf(
		"[chat] ready_to_compile=%t unknowns=%v questions=%v",
		clarified.ReadyToCompile,
		clarified.Unknowns,
		clarified.QuestionsForUser,
	)

	messages, err := h.store.AppendChat(projectID, message, agentReply)
	if err != nil {
		h.notFound(c)
		return
	}

	log.Printf("[chat] reply=%q total_messages=%d", agentReply, len(messages))
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *Handler) buildAgentReply(message string) string {
	parsed, err := h.compiler.Parse(message)
	if err != nil {
		return "我已收到你的补充信息，接下来会继续确认需求边界并推进 IR 生成。"
	}

	intent := strings.TrimSpace(parsed.Intent)
	if intent == "" {
		intent = "需求澄清"
	}

	taskPreview := "继续梳理执行步骤"
	if len(parsed.Tasks) > 0 {
		taskPreview = parsed.Tasks[0]
	}

	constraintPreview := "保留关键风险控制"
	if len(parsed.Constraints) > 0 {
		constraintPreview = parsed.Constraints[0]
	}

	goalPreview := strings.TrimSpace(parsed.Goal)
	if goalPreview == "" {
		goalPreview = message
	}
	if len([]rune(goalPreview)) > 42 {
		goalPreview = string([]rune(goalPreview)[:42]) + "..."
	}

	return "收到，这轮我识别到的重点目标是“" + goalPreview + "”。" +
		"当前意图是“" + intent + "”，我会先推进“" + taskPreview + "”，并遵循“" + constraintPreview + "”这条约束。"
}

func (h *Handler) getChat(c *gin.Context) {
	messages, err := h.store.GetConversation(c.Param("id"))
	if err != nil {
		h.notFound(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *Handler) compile(c *gin.Context) {
	projectID := c.Param("id")

	requirement, err := h.store.GetLatestRequirement(projectID)
	if err != nil {
		h.notFound(c)
		return
	}

	history, err := h.store.GetConversation(projectID)
	if err != nil {
		h.notFound(c)
		return
	}

	clarified, err := h.compiler.Clarify(toClarifyInput(projectID, history, requirement))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "需求澄清失败"})
		return
	}

	parsed, err := h.compiler.ParseFromClarification(clarified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "需求编译失败"})
		return
	}

	ir := model.IRArtifact{
		Version:             parsed.Version,
		Intent:              parsed.Intent,
		Goal:                parsed.Goal,
		Entities:            parsed.Entities,
		Tasks:               parsed.Tasks,
		Constraints:         parsed.Constraints,
		Tools:               parsed.Tools,
		AcceptanceChecklist: parsed.AcceptanceChecklist,
		CompletionCriteria:  parsed.CompletionCriteria,
		HumanGates:          parsed.HumanGates,
		GeneratedAt:         firstNonEmpty(parsed.GeneratedAt, time.Now().Format("2006-01-02 15:04:05")),
	}

	saved, err := h.store.SaveIR(projectID, ir)
	if err != nil {
		h.notFound(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ir": saved})
}

func (h *Handler) assemble(c *gin.Context) {
	graph, err := h.store.Assemble(c.Param("id"))
	if err != nil {
		if err == store.ErrProjectNotFound {
			h.notFound(c)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agentGraph": graph})
}

func (h *Handler) artifacts(c *gin.Context) {
	artifacts, err := h.store.GetArtifacts(c.Param("id"))
	if err != nil {
		h.notFound(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"artifacts": artifacts})
}

func (h *Handler) notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
}

func toClarifyInput(projectID string, history []model.ChatMessage, latest string) compiler.ClarifyInput {
	items := make([]compiler.ClarifyHistoryMessage, 0, len(history))
	for _, message := range history {
		items = append(items, compiler.ClarifyHistoryMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return compiler.ClarifyInput{
		ProjectID:         projectID,
		History:           items,
		LatestUserMessage: latest,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
