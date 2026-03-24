package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ment-agents/backend/internal/app"
	"ment-agents/backend/internal/compiler"
	"ment-agents/backend/internal/config"
	"ment-agents/backend/internal/promptassets"
	"ment-agents/backend/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = "root:root@tcp(127.0.0.1:3306)/ment_agents?charset=utf8mb4&parseTime=True&loc=Local"
	}

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	promptBundle, err := promptassets.Load(cfg.PromptAssets)
	if err != nil {
		log.Fatalf("load prompt assets failed: %v", err)
	}

	st, err := store.NewStore(mysqlDSN)
	if err != nil {
		log.Fatalf("init store failed: %v", err)
	}

	compilerSvc, err := compiler.NewService(cfg.Compiler, promptBundle)
	if err != nil {
		log.Fatalf("init compiler failed: %v", err)
	}

	r := gin.Default()
	h := app.NewHandler(st, compilerSvc)
	h.RegisterRoutes(r)

	log.Printf("backend listening on :%s", port)
	log.Printf("compiler mode: %s", cfg.Compiler.Mode)
	log.Printf("prompt assets loaded: %d stages", promptBundle.StageCount())

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
