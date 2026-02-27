package handler

import (
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/logger"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InterventionHandler struct {
	service service.InterventionService
}

func NewInterventionHandler(service service.InterventionService) *InterventionHandler {
	return &InterventionHandler{service: service}
}

// Query 查询干预库
func (h *InterventionHandler) Query(c *gin.Context) {
	var filters domain.InterventionQueryDTO
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
		logger.Error("查询干预库失败", zap.Error(err))
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

// Add 添加干预库记录
func (h *InterventionHandler) Add(c *gin.Context) {
	var req domain.InterventionAddDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	intervention, err := h.service.Add(c.Request.Context(), &req)
	if err != nil {
		logger.Error("添加干预库记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   intervention,
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Edit 编辑干预库记录
func (h *InterventionHandler) Edit(c *gin.Context) {
	var req domain.InterventionEditDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "参数错误",
			ErrNo:  400,
		})
		return
	}

	intervention, err := h.service.Edit(c.Request.Context(), req.ID, &req)
	if err != nil {
		logger.Error("编辑干预库记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			Data:   nil,
			ErrMsg: "服务暂时不可用",
			ErrNo:  500,
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   intervention,
		ErrMsg: "",
		ErrNo:  0,
	})
}

// Delete 删除干预库记录
func (h *InterventionHandler) Delete(c *gin.Context) {
	var req domain.InterventionDeleteDTO
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
		logger.Error("删除干预库记录失败", zap.Error(err))
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
