package promptassets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ment-agents/backend/internal/config"
)

const (
	StageClarification = "clarification"
	StageCompile       = "compile"
	StageAssemble      = "assemble"
	StageOrchestrate   = "orchestrate"
	StageReflection    = "reflection"
	StageDelivery      = "delivery"
)

type StagePrompt struct {
	System    string
	Validator string
}

type Bundle struct {
	RootDir      string
	ManifestPath string
	BaseSystem   string
	Stages       map[string]StagePrompt
}

func Load(cfg config.PromptAssetsConfig) (*Bundle, error) {
	rootDir := filepath.Clean(cfg.RootDir)
	manifestPath := filepath.Clean(cfg.Manifest)

	if _, err := os.Stat(manifestPath); err != nil {
		return nil, fmt.Errorf("prompt manifest unavailable: %w", err)
	}

	baseSystem, err := readPromptFile(rootDir, cfg.BaseSystem)
	if err != nil {
		return nil, err
	}

	stages := map[string]config.PromptFilePair{
		StageClarification: cfg.Clarification,
		StageCompile:       cfg.Compile,
		StageAssemble:      cfg.Assemble,
		StageOrchestrate:   cfg.Orchestrate,
		StageReflection:    cfg.Reflection,
		StageDelivery:      cfg.Delivery,
	}

	bundle := &Bundle{
		RootDir:      rootDir,
		ManifestPath: manifestPath,
		BaseSystem:   baseSystem,
		Stages:       make(map[string]StagePrompt, len(stages)),
	}

	for stage, pair := range stages {
		systemPrompt, err := readPromptFile(rootDir, pair.System)
		if err != nil {
			return nil, fmt.Errorf("load %s system prompt failed: %w", stage, err)
		}
		validatorPrompt, err := readPromptFile(rootDir, pair.Validator)
		if err != nil {
			return nil, fmt.Errorf("load %s validator prompt failed: %w", stage, err)
		}
		bundle.Stages[stage] = StagePrompt{
			System:    systemPrompt,
			Validator: validatorPrompt,
		}
	}

	return bundle, nil
}

func (b *Bundle) StageCount() int {
	if b == nil {
		return 0
	}
	return len(b.Stages)
}

func (b *Bundle) SystemPrompt(stage string) string {
	if b == nil {
		return ""
	}
	stagePrompt, ok := b.Stages[stage]
	if !ok {
		return strings.TrimSpace(b.BaseSystem)
	}
	return composePrompt(b.BaseSystem, stagePrompt.System)
}

func (b *Bundle) ValidatorPrompt(stage string) string {
	if b == nil {
		return ""
	}
	stagePrompt, ok := b.Stages[stage]
	if !ok {
		return ""
	}
	return strings.TrimSpace(stagePrompt.Validator)
}

func composePrompt(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, "\n\n")
}

func readPromptFile(rootDir, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", nil
	}

	path := name
	if !filepath.IsAbs(path) {
		path = filepath.Join(rootDir, name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read prompt file %s failed: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}
