package handler

import (
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/logger"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RiskEngineHandler struct {
	service service.RiskEngineService
}

func NewRiskEngineHandler(service service.RiskEngineService) *RiskEngineHandler {
	return &RiskEngineHandler{service: service}
}

// Query 查询防护策略
func (h *RiskEngineHandler) Query(c *gin.Context) {
	var filters domain.RiskEngineQueryDTO
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
		logger.Error("查询防护策略失败", zap.Error(err))
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

// Add 添加防护策略记录
func (h *RiskEngineHandler) Add(c *gin.Context) {
	var req domain.RiskEngineAddDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	engine, err := h.service.Add(c.Request.Context(), &req)
	if err != nil {
		logger.Error("添加防护策略记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   engine,
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Edit 编辑防护策略记录
func (h *RiskEngineHandler) Edit(c *gin.Context) {
	var req domain.RiskEngineEditDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	engine, err := h.service.Edit(c.Request.Context(), req.ID, &req)
	if err != nil {
		logger.Error("编辑防护策略记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   engine,
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Delete 删除防护策略记录
func (h *RiskEngineHandler) Delete(c *gin.Context) {
	var req domain.RiskEngineDeleteDTO
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
		logger.Error("删除防护策略记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   nil,
		ErrMsg: "",
		ErrNo:  0,
	})
}
