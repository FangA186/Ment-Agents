package model

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ChatMessage struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type AcceptanceChecklistItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Required    bool   `json:"required"`
	Status      string `json:"status"`
}

type IRArtifact struct {
	Version             string                    `json:"version"`
	Intent              string                    `json:"intent"`
	Goal                string                    `json:"goal"`
	Entities            map[string]string         `json:"entities"`
	Tasks               []string                  `json:"tasks"`
	Constraints         []string                  `json:"constraints"`
	Tools               []string                  `json:"tools"`
	AcceptanceChecklist []AcceptanceChecklistItem `json:"acceptanceChecklist,omitempty"`
	CompletionCriteria  string                    `json:"completionCriteria,omitempty"`
	HumanGates          []string                  `json:"humanGates,omitempty"`
	GeneratedAt         string                    `json:"generatedAt"`
}

type AgentNode struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	RoleTemplate string   `json:"roleTemplate"`
	Skills       []string `json:"skills"`
	Tools        []string `json:"tools"`
	Memory       []string `json:"memory"`
}

type AgentEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type AgentGraph struct {
	Nodes       []AgentNode `json:"nodes"`
	Edges       []AgentEdge `json:"edges"`
	GeneratedAt string      `json:"generatedAt"`
}

type ProjectArtifacts struct {
	IR         *IRArtifact `json:"ir,omitempty"`
	AgentGraph *AgentGraph `json:"agentGraph,omitempty"`
}
