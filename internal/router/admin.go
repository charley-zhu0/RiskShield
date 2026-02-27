package router

import "github.com/gin-gonic/gin"

func setupAdminRoutes(r *gin.Engine, cfg *Config) {
	admin := r.Group("/risk_shield/admin")
	{
		// 干预库管理路由
		intervention := admin.Group("/intervention")
		{
			intervention.GET("/query", cfg.InterventionHandler.Query)
		}

		adminIntervention := admin.Group("/admin/intervention")
		{
			adminIntervention.POST("/add", cfg.InterventionHandler.Add)
			adminIntervention.POST("/edit", cfg.InterventionHandler.Edit)
			adminIntervention.POST("/delete", cfg.InterventionHandler.Delete)
		}

		// 标签管理路由
		label := admin.Group("/label")
		{
			label.GET("/query", cfg.LabelHandler.Query)
		}

		adminLabel := admin.Group("/admin/label")
		{
			adminLabel.POST("/add", cfg.LabelHandler.Add)
			adminLabel.POST("/edit", cfg.LabelHandler.Edit)
			adminLabel.POST("/delete", cfg.LabelHandler.Delete)
		}

		// 防护策略管理路由
		riskengine := admin.Group("/riskengine")
		{
			riskengine.GET("/query", cfg.RiskEngineHandler.Query)
		}

		adminRiskEngine := admin.Group("/admin/riskengine")
		{
			adminRiskEngine.POST("/add", cfg.RiskEngineHandler.Add)
			adminRiskEngine.POST("/edit", cfg.RiskEngineHandler.Edit)
			adminRiskEngine.POST("/delete", cfg.RiskEngineHandler.Delete)
		}

		// 敏感词管理路由
		sensitiveword := admin.Group("/sensitiveword")
		{
			sensitiveword.GET("/query", cfg.SensitiveWordHandler.Query)
		}

		adminSensitiveWord := admin.Group("/admin/sensitiveword")
		{
			adminSensitiveWord.POST("/adds", cfg.SensitiveWordHandler.Add)
			adminSensitiveWord.POST("/edit", cfg.SensitiveWordHandler.Edit)
			adminSensitiveWord.POST("/delete", cfg.SensitiveWordHandler.Delete)
		}
	}
}
