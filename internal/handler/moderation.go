package handler

import (
	"log"
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	result, err := h.service.Moderate(c.Request.Context(), req.Content)
	if err != nil {
		log.Printf("[ERROR] moderation failed: trace_id=%s, error=%v", req.TraceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务暂时不可用"})
		return
	}

	c.JSON(http.StatusOK, result)
}
