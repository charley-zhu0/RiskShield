package handler

import (
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/logger"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SensitiveWordHandler struct {
	service service.SensitiveWordService
}

func NewSensitiveWordHandler(service service.SensitiveWordService) *SensitiveWordHandler {
	return &SensitiveWordHandler{service: service}
}

// Query 查询敏感词
func (h *SensitiveWordHandler) Query(c *gin.Context) {
	var filters domain.SensitiveWordQueryDTO
	if err := c.ShouldBindQuery(&filters); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	response, err := h.service.Query(c.Request.Context(), &filters)
	if err != nil {
		logger.Error("查询敏感词失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   response,
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Add 批量添加敏感词
func (h *SensitiveWordHandler) Add(c *gin.Context) {
	var req domain.SensitiveWordAddDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	words, err := h.service.Add(c.Request.Context(), &req)
	if err != nil {
		logger.Error("添加敏感词失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data: map[string]interface{}{
			"list":  words,
			"total": len(words),
		},
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Edit 编辑敏感词
func (h *SensitiveWordHandler) Edit(c *gin.Context) {
	var req domain.SensitiveWordEditDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	word, err := h.service.Edit(c.Request.Context(), req.ID, &req)
	if err != nil {
		logger.Error("编辑敏感词失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data: map[string]interface{}{
			"list":  []interface{}{word},
			"total": 1,
		},
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Delete 删除敏感词
func (h *SensitiveWordHandler) Delete(c *gin.Context) {
	var req domain.SensitiveWordDeleteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.ID); err != nil {
		logger.Error("删除敏感词失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data: map[string]interface{}{
			"list":  []interface{}{},
			"total": 0,
		},
		ErrMsg: "",
		ErrNo:  0,
	})
}
