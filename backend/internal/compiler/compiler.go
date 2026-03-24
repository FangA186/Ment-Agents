package compiler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"ment-agents/backend/internal/config"
	"ment-agents/backend/internal/model"
	"ment-agents/backend/internal/promptassets"
)

const (
	defaultClarificationPrompt = `你是需求澄清阶段的分析Agent。
你的职责是把用户输入与历史对话整理成结构化澄清结果，并给出一段自然语言回复。
你必须只返回 JSON，对象字段固定为：
- goal
- businessContext
- constraints
- unknowns
- acceptanceHints
- readyToCompile
- questionsForUser
- userReply
当信息还不完整时，只追问最关键的问题，最多 3 个。`

	defaultCompilePrompt = `你是 Compiler。
你的职责是把已经澄清的需求快照转换成结构化 IR。
你必须只返回 JSON，对象字段固定为：
- version
- intent
- goal
- entities
- tasks
- constraints
- tools
- acceptanceChecklist
- completionCriteria
- humanGates
- generatedAt
不要输出解释文字。`
)

type ParseResult struct {
	Version             string                          `json:"version"`
	Intent              string                          `json:"intent"`
	Goal                string                          `json:"goal"`
	Entities            map[string]string               `json:"entities"`
	Tasks               []string                        `json:"tasks"`
	Constraints         []string                        `json:"constraints"`
	Tools               []string                        `json:"tools"`
	AcceptanceChecklist []model.AcceptanceChecklistItem `json:"acceptanceChecklist,omitempty"`
	CompletionCriteria  string                          `json:"completionCriteria,omitempty"`
	HumanGates          []string                        `json:"humanGates,omitempty"`
	GeneratedAt         string                          `json:"generatedAt,omitempty"`
}

type ClarifyHistoryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClarifyInput struct {
	ProjectID         string                  `json:"projectId"`
	History           []ClarifyHistoryMessage `json:"history"`
	LatestUserMessage string                  `json:"latestUserMessage"`
}

type ClarifyResult struct {
	Goal             string   `json:"goal"`
	BusinessContext  string   `json:"businessContext"`
	Constraints      []string `json:"constraints"`
	Unknowns         []string `json:"unknowns"`
	AcceptanceHints  []string `json:"acceptanceHints"`
	ReadyToCompile   bool     `json:"readyToCompile"`
	QuestionsForUser []string `json:"questionsForUser"`
	UserReply        string   `json:"userReply"`
}

type ValidationResult struct {
	Passed         bool     `json:"passed"`
	FailureReasons []string `json:"failureReasons"`
	ReworkAdvice   string   `json:"reworkAdvice"`
}

type Service struct {
	provider Provider
	timeout  time.Duration
	prompts  *promptassets.Bundle
}

type Provider interface {
	Structured(ctx context.Context, systemPrompt string, userInput any, out any) error
}

func NewService(cfg config.CompilerConfig, prompts *promptassets.Bundle) (*Service, error) {
	var provider Provider

	switch cfg.Mode {
	case config.CompilerModeLocalAPI:
		provider = &LocalAPIProvider{
			Endpoint: cfg.LocalAPI.Endpoint,
			APIKey:   cfg.LocalAPI.APIKey,
			Client:   &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second},
		}
	case config.CompilerModeModelAPI:
		provider = &ModelAPIProvider{
			Endpoint: cfg.ModelAPI.Endpoint,
			APIKey:   cfg.ModelAPI.APIKey,
			Model:    cfg.ModelAPI.Model,
			Client:   &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second},
		}
	default:
		return nil, fmt.Errorf("unsupported compiler mode: %s", cfg.Mode)
	}

	return &Service{
		provider: provider,
		timeout:  time.Duration(cfg.TimeoutSeconds) * time.Second,
		prompts:  prompts,
	}, nil
}

func (s *Service) Clarify(input ClarifyInput) (ClarifyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	systemPrompt := defaultClarificationPrompt
	if s.prompts != nil {
		if prompt := strings.TrimSpace(s.prompts.SystemPrompt(promptassets.StageClarification)); prompt != "" {
			systemPrompt = prompt
		}
	}

	fallback := normalizeClarifyResult(heuristicClarify(input), input.LatestUserMessage)
	var result ClarifyResult
	if err := s.provider.Structured(ctx, systemPrompt, input, &result); err != nil {
		return fallback, nil
	}
	result = normalizeClarifyResult(result, input.LatestUserMessage)
	if !s.validateOutput(ctx, promptassets.StageClarification, result) {
		return fallback, nil
	}
	return result, nil
}

func (s *Service) Parse(requirement string) (ParseResult, error) {
	return s.ParseFromClarification(heuristicClarify(ClarifyInput{LatestUserMessage: requirement}))
}

func (s *Service) ParseFromClarification(snapshot ClarifyResult) (ParseResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	systemPrompt := defaultCompilePrompt
	if s.prompts != nil {
		if prompt := strings.TrimSpace(s.prompts.SystemPrompt(promptassets.StageCompile)); prompt != "" {
			systemPrompt = prompt
		}
	}

	fallback := normalizeParseResult(heuristicParse(snapshot.Goal), snapshot)
	var result ParseResult
	if err := s.provider.Structured(ctx, systemPrompt, snapshot, &result); err != nil {
		return fallback, nil
	}
	result = normalizeParseResult(result, snapshot)
	if !s.validateOutput(ctx, promptassets.StageCompile, result) {
		return fallback, nil
	}
	return result, nil
}

func (s *Service) validateOutput(ctx context.Context, stage string, payload any) bool {
	if s.prompts == nil {
		return true
	}

	validatorPrompt := strings.TrimSpace(s.prompts.ValidatorPrompt(stage))
	if validatorPrompt == "" {
		return true
	}

	var verdict ValidationResult
	if err := s.provider.Structured(ctx, validatorPrompt, payload, &verdict); err != nil {
		return true
	}

	return verdict.Passed
}

type LocalAPIProvider struct {
	Endpoint string
	APIKey   string
	Client   *http.Client
}

func (p *LocalAPIProvider) Structured(ctx context.Context, systemPrompt string, userInput any, out any) error {
	parseOut, ok := out.(*ParseResult)
	if !ok {
		return fmt.Errorf("local api provider only supports parse result output")
	}

	requirement := extractRequirement(userInput)
	payload := map[string]string{
		"requirement":   requirement,
		"system_prompt": systemPrompt,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("local api status: %d", resp.StatusCode)
	}

	var parsed ParseResult
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}
	*parseOut = parsed
	return nil
}

type ModelAPIProvider struct {
	Endpoint string
	APIKey   string
	Model    string
	Client   *http.Client
}

func (p *ModelAPIProvider) Structured(ctx context.Context, systemPrompt string, userInput any, out any) error {
	type chatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type request struct {
		Model       string        `json:"model"`
		Temperature float64       `json:"temperature"`
		Messages    []chatMessage `json:"messages"`
	}
	type response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	body, _ := json.Marshal(request{
		Model:       p.Model,
		Temperature: 0.1,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: marshalUserInput(userInput)},
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("model api status: %d", resp.StatusCode)
	}

	var decoded response
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return err
	}
	if len(decoded.Choices) == 0 {
		return fmt.Errorf("empty model response")
	}

	content := strings.TrimSpace(decoded.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), out); err != nil {
		return fmt.Errorf("unmarshal model response failed: %w", err)
	}
	return nil
}

func normalizeClarifyResult(result ClarifyResult, latest string) ClarifyResult {
	fallback := heuristicClarify(ClarifyInput{LatestUserMessage: latest})

	if strings.TrimSpace(result.Goal) == "" {
		result.Goal = fallback.Goal
	}
	if strings.TrimSpace(result.BusinessContext) == "" {
		result.BusinessContext = fallback.BusinessContext
	}
	if len(result.Constraints) == 0 {
		result.Constraints = fallback.Constraints
	} else {
		result.Constraints = uniqueStrings(result.Constraints)
	}
	if len(result.Unknowns) == 0 && !result.ReadyToCompile {
		result.Unknowns = fallback.Unknowns
	} else {
		result.Unknowns = uniqueStrings(result.Unknowns)
	}
	if len(result.AcceptanceHints) == 0 {
		result.AcceptanceHints = fallback.AcceptanceHints
	} else {
		result.AcceptanceHints = uniqueStrings(result.AcceptanceHints)
	}
	if len(result.QuestionsForUser) == 0 && !result.ReadyToCompile {
		result.QuestionsForUser = fallback.QuestionsForUser
	} else {
		result.QuestionsForUser = uniqueStrings(result.QuestionsForUser)
	}
	if len(result.Unknowns) > 0 {
		result.ReadyToCompile = false
	}
	if strings.TrimSpace(result.UserReply) == "" {
		result.UserReply = buildClarificationReply(result)
	}
	return result
}

func normalizeParseResult(result ParseResult, snapshot ClarifyResult) ParseResult {
	base := heuristicParse(snapshot.Goal)

	if strings.TrimSpace(result.Version) == "" {
		result.Version = "v1"
	}
	if strings.TrimSpace(result.Intent) == "" {
		result.Intent = base.Intent
	}
	if strings.TrimSpace(result.Goal) == "" {
		result.Goal = firstNonEmpty(strings.TrimSpace(snapshot.Goal), base.Goal)
	}
	if result.Entities == nil {
		result.Entities = base.Entities
	}
	if len(result.Tasks) == 0 {
		result.Tasks = base.Tasks
	}
	if len(result.Constraints) == 0 {
		result.Constraints = uniqueStrings(append(snapshot.Constraints, base.Constraints...))
	} else {
		result.Constraints = uniqueStrings(result.Constraints)
	}
	if len(result.Tools) == 0 {
		result.Tools = base.Tools
	}
	if len(result.AcceptanceChecklist) == 0 {
		result.AcceptanceChecklist = buildAcceptanceChecklist(result.Goal, result.Tasks, snapshot.AcceptanceHints)
	}
	if strings.TrimSpace(result.CompletionCriteria) == "" {
		result.CompletionCriteria = "验收清单全部通过，且用户完成最终确认。"
	}
	if len(result.HumanGates) == 0 {
		result.HumanGates = []string{"ir_confirm", "final_confirm"}
	} else {
		result.HumanGates = uniqueStrings(result.HumanGates)
	}
	if strings.TrimSpace(result.GeneratedAt) == "" {
		result.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	result.Tasks = uniqueStrings(result.Tasks)
	result.Tools = uniqueStrings(result.Tools)
	return result
}

func heuristicClarify(input ClarifyInput) ClarifyResult {
	latest := strings.TrimSpace(input.LatestUserMessage)
	combined := buildUserContext(input.History, latest)

	switch {
	case containsAny(latest, "你好", "您好", "hello", "hi"):
		return ClarifyResult{
			Goal:             "待用户提供具体业务需求",
			BusinessContext:  "当前仍处于需求澄清起点。",
			Constraints:      []string{},
			Unknowns:         []string{"尚未提供具体业务目标", "尚未提供约束条件"},
			AcceptanceHints:  []string{"确认业务目标", "确认约束边界"},
			ReadyToCompile:   false,
			QuestionsForUser: []string{"你希望搭建什么系统或自动化流程？"},
			UserReply:        "我在这里。你可以直接告诉我想实现的业务目标、关键约束和期望结果，我会先帮你把需求澄清清楚。",
		}
	case containsAny(latest, "你是谁", "你是?", "介绍一下你自己", "你是做什么的"):
		return ClarifyResult{
			Goal:             "待用户提供具体业务需求",
			BusinessContext:  "当前为身份说明场景。",
			Constraints:      []string{},
			Unknowns:         []string{"尚未提供具体业务需求"},
			AcceptanceHints:  []string{"补充业务目标"},
			ReadyToCompile:   false,
			QuestionsForUser: []string{"你希望我帮你梳理什么业务需求？"},
			UserReply:        "我是负责需求澄清的分析Agent，会先把你的目标、约束和验收标准梳理成结构化结果，再推进 IR 生成。你可以直接告诉我想做什么。",
		}
	}

	goal := deriveGoal(combined)
	constraints := deriveConstraints(combined)
	unknowns := deriveUnknowns(combined)
	acceptanceHints := deriveAcceptanceHints(combined)
	ready := len(unknowns) == 0 && goal != ""
	questions := deriveClarificationQuestions(unknowns)

	return ClarifyResult{
		Goal:             goal,
		BusinessContext:  deriveBusinessContext(combined),
		Constraints:      constraints,
		Unknowns:         unknowns,
		AcceptanceHints:  acceptanceHints,
		ReadyToCompile:   ready,
		QuestionsForUser: questions,
		UserReply:        buildClarificationReply(ClarifyResult{Goal: goal, Constraints: constraints, Unknowns: unknowns, QuestionsForUser: questions, ReadyToCompile: ready}),
	}
}

func heuristicParse(requirement string) ParseResult {
	req := strings.TrimSpace(requirement)
	entities := map[string]string{}

	if strings.Contains(req, "预算") {
		re := regexp.MustCompile(`预算[^0-9]*(\d+[kKwW万]?)`)
		if match := re.FindStringSubmatch(req); len(match) == 2 {
			entities["budget"] = match[1]
		}
	}
	if strings.Contains(req, "人工") && strings.Contains(req, "审核") {
		entities["human_review"] = "required"
	}
	if containsAny(req, "日报", "周报", "报告") {
		entities["reporting"] = "enabled"
	}

	return ParseResult{
		Version:  "v1",
		Intent:   inferIntent(req),
		Goal:     firstNonEmpty(req, "待补充具体业务目标"),
		Entities: entities,
		Tasks: []string{
			"澄清业务目标与范围",
			"拆解执行步骤并分配角色",
			"输出阶段结果与回溯信息",
		},
		Constraints: []string{
			"保留人工审核入口",
			"关键步骤支持失败回滚",
		},
		Tools: []string{
			"HTTP Client",
			"SQL Query",
			"日志写入",
		},
		AcceptanceChecklist: []model.AcceptanceChecklistItem{
			{
				ID:          "ac-001",
				Title:       "核心目标已结构化",
				Description: "IR 已覆盖目标、任务、约束与工具。",
				Owner:       "分析Agent",
				Required:    true,
				Status:      "pending",
			},
		},
		CompletionCriteria: "验收清单全部通过，且用户完成最终确认。",
		HumanGates:         []string{"ir_confirm", "final_confirm"},
		GeneratedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}
}

func deriveGoal(text string) string {
	text = strings.TrimSpace(strings.TrimPrefix(text, "需求:"))
	if text == "" {
		return ""
	}
	if len([]rune(text)) > 60 {
		return string([]rune(text)[:60]) + "..."
	}
	return text
}

func deriveBusinessContext(text string) string {
	switch {
	case containsAny(text, "营销", "投放", "活动"):
		return "该需求聚焦营销活动自动化、执行协同与结果复盘。"
	case containsAny(text, "工单", "售后", "客服"):
		return "该需求聚焦售后工单流转、分派与处理闭环。"
	case containsAny(text, "内容", "分发", "增长"):
		return "该需求聚焦内容生产、分发协同与增长反馈。"
	default:
		return "该需求仍处于通用业务流程梳理阶段。"
	}
}

func deriveConstraints(text string) []string {
	constraints := []string{}
	if strings.Contains(text, "预算") {
		re := regexp.MustCompile(`预算[^0-9]*(\d+[kKwW万]?)`)
		if match := re.FindStringSubmatch(text); len(match) == 2 {
			constraints = append(constraints, "预算控制在 "+match[1])
		}
	}
	if containsAny(text, "人工审核", "人工审批", "人工介入") {
		constraints = append(constraints, "保留人工审核入口")
	}
	if containsAny(text, "可回滚", "回滚") {
		constraints = append(constraints, "关键步骤支持失败回滚")
	}
	if containsAny(text, "追踪", "回溯", "日志") {
		constraints = append(constraints, "执行链路需可追踪")
	}
	return uniqueStrings(constraints)
}

func deriveUnknowns(text string) []string {
	unknowns := []string{}
	if strings.TrimSpace(text) == "" {
		return []string{"尚未提供具体业务目标", "尚未提供约束条件"}
	}
	if !containsAny(text, "预算", "成本", "额度") {
		unknowns = append(unknowns, "预算范围未明确")
	}
	if !containsAny(text, "人工审核", "人工审批", "人工介入", "自动通过") {
		unknowns = append(unknowns, "是否需要人工审核未明确")
	}
	if !containsAny(text, "日报", "周报", "报告", "输出", "交付", "验收") {
		unknowns = append(unknowns, "最终交付物或验收口径未明确")
	}
	return uniqueStrings(unknowns)
}

func deriveAcceptanceHints(text string) []string {
	hints := []string{"可生成结构化 IR"}
	if containsAny(text, "日报", "周报", "报告") {
		hints = append(hints, "需要输出报告或复盘结果")
	}
	if containsAny(text, "人工审核", "人工审批", "人工介入") {
		hints = append(hints, "需保留人工闸门")
	}
	if containsAny(text, "追踪", "回溯", "日志") {
		hints = append(hints, "需保留执行过程追踪信息")
	}
	return uniqueStrings(hints)
}

func deriveClarificationQuestions(unknowns []string) []string {
	questions := make([]string, 0, 3)
	for _, unknown := range unknowns {
		switch unknown {
		case "预算范围未明确":
			questions = append(questions, "这次项目的预算范围大概是多少？")
		case "是否需要人工审核未明确":
			questions = append(questions, "是否需要保留人工审核或人工审批入口？")
		case "最终交付物或验收口径未明确":
			questions = append(questions, "你希望最终交付物是什么形式，例如日报、看板还是执行清单？")
		default:
			questions = append(questions, "请补充当前缺失的关键信息。")
		}
		if len(questions) == 3 {
			break
		}
	}
	return uniqueStrings(questions)
}

func buildClarificationReply(result ClarifyResult) string {
	if result.ReadyToCompile {
		return fmt.Sprintf("我已经把这轮需求收敛得比较完整了。当前目标是“%s”，接下来我会进入 IR 结构化生成；如果你没有补充项，我们就可以继续。", firstNonEmpty(result.Goal, "当前需求"))
	}

	parts := []string{}
	if result.Goal != "" {
		parts = append(parts, "我先整理到的目标是“"+result.Goal+"”。")
	}
	if len(result.Constraints) > 0 {
		parts = append(parts, "已识别约束："+strings.Join(result.Constraints, "；")+"。")
	}
	if len(result.QuestionsForUser) > 0 {
		parts = append(parts, "为了继续推进，我还需要确认："+strings.Join(result.QuestionsForUser, " "))
	} else if len(result.Unknowns) > 0 {
		parts = append(parts, "当前还有一些关键信息未明确，我建议先把缺失项补齐。")
	}
	if len(parts) == 0 {
		return "我先帮你把需求澄清清楚，再推进 IR 生成。请继续补充你的目标、约束和期望结果。"
	}
	return strings.Join(parts, " ")
}

func buildAcceptanceChecklist(goal string, tasks, hints []string) []model.AcceptanceChecklistItem {
	items := []model.AcceptanceChecklistItem{
		{
			ID:          "ac-001",
			Title:       "目标定义清晰",
			Description: "IR 中已完整记录业务目标与关键约束。",
			Owner:       "分析Agent",
			Required:    true,
			Status:      "pending",
		},
	}

	if len(tasks) > 0 {
		items = append(items, model.AcceptanceChecklistItem{
			ID:          "ac-002",
			Title:       "任务拆解完整",
			Description: "核心任务已经拆解并可供后续 Agent 执行。",
			Owner:       "规划Agent",
			Required:    true,
			Status:      "pending",
		})
	}
	if len(hints) > 0 || goal != "" {
		items = append(items, model.AcceptanceChecklistItem{
			ID:          "ac-003",
			Title:       "交付口径可验收",
			Description: "存在明确的交付预期或验收提示，可支撑后续 Reflection 检查。",
			Owner:       "评估Agent",
			Required:    true,
			Status:      "pending",
		})
	}
	return items
}

func marshalUserInput(userInput any) string {
	switch v := userInput.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		raw, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", userInput)
		}
		return string(raw)
	}
}

func extractRequirement(userInput any) string {
	switch v := userInput.(type) {
	case string:
		return strings.TrimSpace(v)
	case ClarifyInput:
		return strings.TrimSpace(v.LatestUserMessage)
	case *ClarifyInput:
		if v == nil {
			return ""
		}
		return strings.TrimSpace(v.LatestUserMessage)
	case ClarifyResult:
		return firstNonEmpty(strings.TrimSpace(v.Goal), strings.Join(v.AcceptanceHints, "；"))
	case *ClarifyResult:
		if v == nil {
			return ""
		}
		return firstNonEmpty(strings.TrimSpace(v.Goal), strings.Join(v.AcceptanceHints, "；"))
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", userInput))
	}
}

func buildUserContext(history []ClarifyHistoryMessage, latest string) string {
	parts := []string{}
	for _, item := range history {
		if item.Role != "user" {
			continue
		}
		content := strings.TrimSpace(item.Content)
		if content != "" {
			parts = append(parts, content)
		}
	}
	if strings.TrimSpace(latest) != "" {
		parts = append(parts, strings.TrimSpace(latest))
	}
	return strings.Join(uniqueStrings(parts), "；")
}

func inferIntent(text string) string {
	switch {
	case containsAny(text, "你好", "您好", "hello", "hi"):
		return "问候"
	case containsAny(text, "你是谁", "你是?", "介绍一下"):
		return "身份询问"
	case containsAny(text, "营销", "投放", "活动"):
		return "营销自动化需求"
	case containsAny(text, "工单", "售后", "客服"):
		return "工单协同需求"
	case containsAny(text, "内容", "分发", "增长"):
		return "内容增长需求"
	default:
		return "需求解析"
	}
}

func containsAny(text string, keywords ...string) bool {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func uniqueStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
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
