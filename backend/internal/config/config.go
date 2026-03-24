package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type CompilerMode string

const (
	CompilerModeLocalAPI CompilerMode = "local_api"
	CompilerModeModelAPI CompilerMode = "model_api"
)

type Config struct {
	Compiler     CompilerConfig     `yaml:"compiler"`
	PromptAssets PromptAssetsConfig `yaml:"prompt_assets"`
}

type CompilerConfig struct {
	Mode           CompilerMode   `yaml:"mode"`
	TimeoutSeconds int            `yaml:"timeout_seconds"`
	LocalAPI       LocalAPIConfig `yaml:"local_api"`
	ModelAPI       ModelAPIConfig `yaml:"model_api"`
}

type LocalAPIConfig struct {
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
}

type ModelAPIConfig struct {
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
}

type PromptAssetsConfig struct {
	RootDir       string         `yaml:"root_dir"`
	Manifest      string         `yaml:"manifest"`
	BaseSystem    string         `yaml:"base_system"`
	Clarification PromptFilePair `yaml:"clarification"`
	Compile       PromptFilePair `yaml:"compile"`
	Assemble      PromptFilePair `yaml:"assemble"`
	Orchestrate   PromptFilePair `yaml:"orchestrate"`
	Reflection    PromptFilePair `yaml:"reflection"`
	Delivery      PromptFilePair `yaml:"delivery"`
}

type PromptFilePair struct {
	System    string `yaml:"system"`
	Validator string `yaml:"validator"`
}

func defaultConfig() Config {
	return Config{
		Compiler: CompilerConfig{
			Mode:           CompilerModeLocalAPI,
			TimeoutSeconds: 20,
			LocalAPI: LocalAPIConfig{
				Endpoint: "http://127.0.0.1:9090/parse",
			},
			ModelAPI: ModelAPIConfig{
				Endpoint: "https://api.openai.com/v1/chat/completions",
				Model:    "gpt-4o-mini",
			},
		},
		PromptAssets: PromptAssetsConfig{
			RootDir:    "./config/prompts",
			Manifest:   "./config/prompts/manifest.yaml",
			BaseSystem: "00-base-system.md",
			Clarification: PromptFilePair{
				System:    "01-clarification-system.md",
				Validator: "01-clarification-validator.md",
			},
			Compile: PromptFilePair{
				System:    "02-compile-system.md",
				Validator: "02-compile-validator.md",
			},
			Assemble: PromptFilePair{
				System:    "03-assemble-system.md",
				Validator: "03-assemble-validator.md",
			},
			Orchestrate: PromptFilePair{
				System:    "04-orchestrate-system.md",
				Validator: "04-orchestrate-validator.md",
			},
			Reflection: PromptFilePair{
				System:    "05-reflection-system.md",
				Validator: "05-reflection-validator.md",
			},
			Delivery: PromptFilePair{
				System:    "06-delivery-system.md",
				Validator: "06-delivery-validator.md",
			},
		},
	}
}

func Load(path string) (Config, error) {
	cfg := defaultConfig()

	if path != "" {
		if raw, err := os.ReadFile(path); err == nil {
			if err := yaml.Unmarshal(raw, &cfg); err != nil {
				return Config{}, fmt.Errorf("parse config yaml failed: %w", err)
			}
		}
	}

	applyEnvOverrides(&cfg)
	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if mode := os.Getenv("COMPILER_MODE"); mode != "" {
		cfg.Compiler.Mode = CompilerMode(mode)
	}
	if endpoint := os.Getenv("COMPILER_LOCAL_API_ENDPOINT"); endpoint != "" {
		cfg.Compiler.LocalAPI.Endpoint = endpoint
	}
	if key := os.Getenv("COMPILER_LOCAL_API_KEY"); key != "" {
		cfg.Compiler.LocalAPI.APIKey = key
	}
	if endpoint := os.Getenv("COMPILER_MODEL_API_ENDPOINT"); endpoint != "" {
		cfg.Compiler.ModelAPI.Endpoint = endpoint
	}
	if key := os.Getenv("COMPILER_MODEL_API_KEY"); key != "" {
		cfg.Compiler.ModelAPI.APIKey = key
	}
	if model := os.Getenv("COMPILER_MODEL_NAME"); model != "" {
		cfg.Compiler.ModelAPI.Model = model
	}
	if timeout := os.Getenv("COMPILER_TIMEOUT_SECONDS"); timeout != "" {
		if sec, err := strconv.Atoi(timeout); err == nil {
			cfg.Compiler.TimeoutSeconds = sec
		}
	}
	if root := os.Getenv("PROMPT_ASSETS_ROOT_DIR"); root != "" {
		cfg.PromptAssets.RootDir = root
	}
	if manifest := os.Getenv("PROMPT_ASSETS_MANIFEST"); manifest != "" {
		cfg.PromptAssets.Manifest = manifest
	}
}

func validate(cfg Config) error {
	if cfg.Compiler.TimeoutSeconds <= 0 {
		return fmt.Errorf("compiler.timeout_seconds must be > 0")
	}
	switch cfg.Compiler.Mode {
	case CompilerModeLocalAPI:
		if cfg.Compiler.LocalAPI.Endpoint == "" {
			return fmt.Errorf("compiler.local_api.endpoint is required")
		}
	case CompilerModeModelAPI:
		if cfg.Compiler.ModelAPI.Endpoint == "" {
			return fmt.Errorf("compiler.model_api.endpoint is required")
		}
		if cfg.Compiler.ModelAPI.Model == "" {
			return fmt.Errorf("compiler.model_api.model is required")
		}
	default:
		return fmt.Errorf("unsupported compiler.mode: %s", cfg.Compiler.Mode)
	}
	if cfg.PromptAssets.RootDir == "" {
		return fmt.Errorf("prompt_assets.root_dir is required")
	}
	if cfg.PromptAssets.Manifest == "" {
		return fmt.Errorf("prompt_assets.manifest is required")
	}
	return nil
}
