package router

import (
	"github.com/charley/riskshield/internal/handler"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ModerationHandler    *handler.ModerationHandler
	InterventionHandler  *handler.InterventionHandler
	LabelHandler         *handler.LabelHandler
	RiskEngineHandler    *handler.RiskEngineHandler
	SensitiveWordHandler *handler.SensitiveWordHandler
}

func Setup(cfg *Config) *gin.Engine {
	r := gin.Default()
	setupPrivateRoutes(r, cfg)
	setupAdminRoutes(r, cfg)
	return r
}
