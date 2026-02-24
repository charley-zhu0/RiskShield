package handler

import (
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/logger"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ModerationHandler struct {
	service service.ModerationService
}

func NewModerationHandler(service service.ModerationService) *ModerationHandler {
	return &ModerationHandler{service: service}
}

func (h *ModerationHandler) Moderate(c *gin.Context) {
	var req domain.ModerationRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	logger.Info("开始审核", zap.String("trace_id", req.TraceID))

	result, err := h.service.Moderate(c.Request.Context(), req.Content)
	if err != nil {
		logger.Error("审核失败", zap.String("trace_id", req.TraceID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务暂时不可用"})
		return
	}

	logger.Info("审核完成", zap.String("trace_id", req.TraceID), zap.String("decision", result.Decision))
	c.JSON(http.StatusOK, result)
}
