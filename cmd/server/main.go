package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charley/riskshield/internal/client"
	"github.com/charley/riskshield/internal/config"
	"github.com/charley/riskshield/internal/handler"
	"github.com/charley/riskshield/internal/logger"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupSignalHandler(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	if err := logger.Init(cfg.LogDir, cfg.LogMaxBackups); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("服务启动中")

	swClient := client.NewSensitiveWordClient(cfg.SensitiveWordURL, cfg.RequestTimeout)
	llmClient := client.NewLLMClient(cfg.LLMServiceURL, cfg.RequestTimeout)
	moderationSvc := service.NewModerationService(swClient, llmClient)
	moderationHandler := handler.NewModerationHandler(moderationSvc)

	r := gin.Default()
	r.POST("/v1/private/moderations/text", moderationHandler.Moderate)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		logger.Info("服务启动在端口", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务启动失败", zap.Error(err))
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	setupSignalHandler(quit)
	<-quit

	logger.Info("收到退出信号，开始优雅关闭")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务关闭失败", zap.Error(err))
		log.Fatalf("服务关闭失败: %v", err)
	}

	logger.Info("服务已安全退出")
}
