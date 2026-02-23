package main

import (
	"log"

	"github.com/charley/riskshield/internal/client"
	"github.com/charley/riskshield/internal/config"
	"github.com/charley/riskshield/internal/handler"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 初始化客户端
	swClient := client.NewSensitiveWordClient(cfg.SensitiveWordURL, cfg.RequestTimeout)
	llmClient := client.NewLLMClient(cfg.LLMServiceURL, cfg.RequestTimeout)

	// 初始化服务
	moderationSvc := service.NewModerationService(swClient, llmClient)

	// 初始化 Handler
	moderationHandler := handler.NewModerationHandler(moderationSvc)

	// 设置路由
	r := gin.Default()
	r.POST("/v1/private/moderations/text", moderationHandler.Moderate)

	// 启动服务
	log.Printf("服务启动在端口 %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
