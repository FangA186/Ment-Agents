package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ment-agents/backend/internal/model"
)

var ErrProjectNotFound = errors.New("project not found")

type projectRecord struct {
	ID        uint      `gorm:"primaryKey"`
	PublicID  string    `gorm:"size:32;uniqueIndex;not null"`
	Name      string    `gorm:"size:255;not null"`
	Status    string    `gorm:"size:32;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (projectRecord) TableName() string { return "projects" }

type chatMessageRecord struct {
	ID        uint      `gorm:"primaryKey"`
	PublicID  string    `gorm:"size:32;uniqueIndex;not null"`
	ProjectID uint      `gorm:"index;not null"`
	Role      string    `gorm:"size:16;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (chatMessageRecord) TableName() string { return "chat_messages" }

type artifactRecord struct {
	ID             uint      `gorm:"primaryKey"`
	ProjectID      uint      `gorm:"uniqueIndex;not null"`
	IRJSON         string    `gorm:"type:longtext"`
	AgentGraphJSON string    `gorm:"type:longtext"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (artifactRecord) TableName() string { return "project_artifacts" }

type Store struct {
	db *gorm.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql failed: %w", err)
	}

	if err := db.AutoMigrate(&projectRecord{}, &chatMessageRecord{}, &artifactRecord{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	s := &Store{db: db}
	if err := s.seedIfEmpty(); err != nil {
		return nil, err
	}

	return s, nil
}

func ts(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func parseProjSeq(publicID string) int {
	re := regexp.MustCompile(`^proj-(\d+)$`)
	m := re.FindStringSubmatch(publicID)
	if len(m) != 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

func (s *Store) nextProjectPublicID(tx *gorm.DB) (string, error) {
	var records []projectRecord
	if err := tx.Select("public_id").Find(&records).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, item := range records {
		if seq := parseProjSeq(item.PublicID); seq > maxSeq {
			maxSeq = seq
		}
	}

	return fmt.Sprintf("proj-%03d", maxSeq+1), nil
}

func (s *Store) nextMessagePublicID(tx *gorm.DB) (string, error) {
	var rec chatMessageRecord
	err := tx.Order("id desc").First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "m-1", nil
		}
		return "", err
	}
	return fmt.Sprintf("m-%d", rec.ID+1), nil
}

func (s *Store) seedIfEmpty() error {
	var count int64
	if err := s.db.Model(&projectRecord{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	project := projectRecord{
		PublicID:  "proj-001",
		Name:      "营销活动自动化系统",
		Status:    "进行中",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.db.Create(&project).Error; err != nil {
		return err
	}

	messages := []chatMessageRecord{
		{PublicID: "m-1", ProjectID: project.ID, Role: "user", Content: "我要构建一个营销活动自动化系统。", CreatedAt: now},
		{PublicID: "m-2", ProjectID: project.ID, Role: "agent", Content: "收到，我会先确认边界条件，再生成 IR。", CreatedAt: now},
	}
	if err := s.db.Create(&messages).Error; err != nil {
		return err
	}

	return nil
}

func toProjectModel(rec projectRecord) model.Project {
	return model.Project{
		ID:        rec.PublicID,
		Name:      rec.Name,
		Status:    rec.Status,
		CreatedAt: ts(rec.CreatedAt),
		UpdatedAt: ts(rec.UpdatedAt),
	}
}

func toMessageModel(rec chatMessageRecord) model.ChatMessage {
	return model.ChatMessage{
		ID:        rec.PublicID,
		Role:      rec.Role,
		Content:   rec.Content,
		CreatedAt: ts(rec.CreatedAt),
	}
}

func (s *Store) ListProjects() []model.Project {
	var records []projectRecord
	_ = s.db.Order("created_at asc").Find(&records).Error

	result := make([]model.Project, 0, len(records))
	for _, rec := range records {
		result = append(result, toProjectModel(rec))
	}
	return result
}

func (s *Store) GetProject(id string) (model.Project, error) {
	var rec projectRecord
	if err := s.db.Where("public_id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Project{}, ErrProjectNotFound
		}
		return model.Project{}, err
	}
	return toProjectModel(rec), nil
}

func (s *Store) CreateProject(name string) model.Project {
	var out model.Project
	_ = s.db.Transaction(func(tx *gorm.DB) error {
		publicID, err := s.nextProjectPublicID(tx)
		if err != nil {
			return err
		}

		rec := projectRecord{PublicID: publicID, Name: name, Status: "规划中"}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
		out = toProjectModel(rec)
		return nil
	})

	return out
}

func (s *Store) AppendChat(projectID, userContent, agentReply string) ([]model.ChatMessage, error) {
	var result []model.ChatMessage

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var project projectRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("public_id = ?", projectID).First(&project).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProjectNotFound
			}
			return err
		}

		userID, err := s.nextMessagePublicID(tx)
		if err != nil {
			return err
		}
		if err := tx.Create(&chatMessageRecord{PublicID: userID, ProjectID: project.ID, Role: "user", Content: userContent}).Error; err != nil {
			return err
		}

		agentID, err := s.nextMessagePublicID(tx)
		if err != nil {
			return err
		}
		if err := tx.Create(&chatMessageRecord{PublicID: agentID, ProjectID: project.ID, Role: "agent", Content: agentReply}).Error; err != nil {
			return err
		}

		if err := tx.Model(&project).Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		var records []chatMessageRecord
		if err := tx.Where("project_id = ?", project.ID).Order("id asc").Find(&records).Error; err != nil {
			return err
		}

		result = make([]model.ChatMessage, 0, len(records))
		for _, item := range records {
			result = append(result, toMessageModel(item))
		}
		return nil
	}); err != nil {
		if err == ErrProjectNotFound {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return result, nil
}

func (s *Store) GetConversation(projectID string) ([]model.ChatMessage, error) {
	var project projectRecord
	if err := s.db.Where("public_id = ?", projectID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	var records []chatMessageRecord
	if err := s.db.Where("project_id = ?", project.ID).Order("id asc").Find(&records).Error; err != nil {
		return nil, err
	}

	result := make([]model.ChatMessage, 0, len(records))
	for _, item := range records {
		result = append(result, toMessageModel(item))
	}
	return result, nil
}

func (s *Store) GetLatestRequirement(projectID string) (string, error) {
	var project projectRecord
	if err := s.db.Where("public_id = ?", projectID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrProjectNotFound
		}
		return "", err
	}

	var latestUser chatMessageRecord
	if err := s.db.Where("project_id = ? AND role = ?", project.ID, "user").Order("id desc").First(&latestUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return project.Name, nil
		}
		return "", err
	}

	if latestUser.Content == "" {
		return project.Name, nil
	}
	return latestUser.Content, nil
}

func (s *Store) SaveIR(projectID string, ir model.IRArtifact) (model.IRArtifact, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var project projectRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("public_id = ?", projectID).First(&project).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProjectNotFound
			}
			return err
		}

		irJSON, err := json.Marshal(ir)
		if err != nil {
			return err
		}

		artifact := artifactRecord{ProjectID: project.ID, IRJSON: string(irJSON)}
		if err := tx.Where("project_id = ?", project.ID).Assign(artifactRecord{IRJSON: string(irJSON)}).FirstOrCreate(&artifact).Error; err != nil {
			return err
		}

		if err := tx.Model(&project).Updates(map[string]any{"status": "进行中", "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == ErrProjectNotFound {
			return model.IRArtifact{}, ErrProjectNotFound
		}
		return model.IRArtifact{}, err
	}

	return ir, nil
}

func (s *Store) Assemble(projectID string) (model.AgentGraph, error) {
	var graph model.AgentGraph
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var project projectRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("public_id = ?", projectID).First(&project).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProjectNotFound
			}
			return err
		}

		var artifact artifactRecord
		if err := tx.Where("project_id = ?", project.ID).First(&artifact).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || artifact.IRJSON == "" {
				return errors.New("ir artifact not found, compile first")
			}
			return err
		}

		graph = model.AgentGraph{
			Nodes: []model.AgentNode{
				{ID: "A1", Name: "分析Agent", RoleTemplate: "需求分析模板", Skills: []string{"需求拆解", "约束提取"}, Tools: []string{"检索工具"}, Memory: []string{"需求上下文"}},
				{ID: "A2", Name: "规划Agent", RoleTemplate: "计划编排模板", Skills: []string{"任务编排", "依赖分析"}, Tools: []string{"任务调度器"}, Memory: []string{"阶段计划"}},
				{ID: "A3", Name: "执行Agent", RoleTemplate: "执行落地模板", Skills: []string{"任务执行", "状态回传"}, Tools: []string{"HTTP 客户端", "数据库查询"}, Memory: []string{"执行日志"}},
			},
			Edges:       []model.AgentEdge{{From: "A1", To: "A2"}, {From: "A2", To: "A3"}},
			GeneratedAt: ts(time.Now()),
		}

		graphJSON, err := json.Marshal(graph)
		if err != nil {
			return err
		}

		if err := tx.Model(&artifact).Updates(map[string]any{"agent_graph_json": string(graphJSON), "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		if err := tx.Model(&project).Update("updated_at", time.Now()).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == ErrProjectNotFound {
			return model.AgentGraph{}, ErrProjectNotFound
		}
		return model.AgentGraph{}, err
	}
	return graph, nil
}

func (s *Store) GetArtifacts(projectID string) (model.ProjectArtifacts, error) {
	var project projectRecord
	if err := s.db.Where("public_id = ?", projectID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ProjectArtifacts{}, ErrProjectNotFound
		}
		return model.ProjectArtifacts{}, err
	}

	var artifact artifactRecord
	if err := s.db.Where("project_id = ?", project.ID).First(&artifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ProjectArtifacts{}, nil
		}
		return model.ProjectArtifacts{}, err
	}

	result := model.ProjectArtifacts{}
	if artifact.IRJSON != "" {
		var ir model.IRArtifact
		if err := json.Unmarshal([]byte(artifact.IRJSON), &ir); err == nil {
			result.IR = &ir
		}
	}
	if artifact.AgentGraphJSON != "" {
		var graph model.AgentGraph
		if err := json.Unmarshal([]byte(artifact.AgentGraphJSON), &graph); err == nil {
			result.AgentGraph = &graph
		}
	}
	return result, nil
}
