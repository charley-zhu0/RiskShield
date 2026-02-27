package router

import "github.com/gin-gonic/gin"

func setupPrivateRoutes(r *gin.Engine, cfg *Config) {
	private := r.Group("/v1/private")
	{
		moderations := private.Group("/moderations")
		{
			moderations.POST("/text", cfg.ModerationHandler.Moderate)
		}
	}
}
