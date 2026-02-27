package handler

import (
	"net/http"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/service"
	"github.com/gin-gonic/gin"
)

type LabelHandler struct {
	labelService service.LabelService
}

// NewLabelHandler 创建标签 Handler 实例
func NewLabelHandler(labelService service.LabelService) *LabelHandler {
	return &LabelHandler{
		labelService: labelService,
	}
}

// Query 查询标签列表
func (h *LabelHandler) Query(c *gin.Context) {
	var dto domain.LabelQueryDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			ErrNo:  400,
			ErrMsg: "参数验证失败: " + err.Error(),
		})
		return
	}

	resp, err := h.labelService.Query(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			ErrNo:  500,
			ErrMsg: "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   resp,
		ErrNo:  0,
		ErrMsg: "",
	})
}

// Add 添加标签
func (h *LabelHandler) Add(c *gin.Context) {
	var dto domain.LabelAddDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			ErrNo:  400,
			ErrMsg: "参数验证失败: " + err.Error(),
		})
		return
	}

	label, err := h.labelService.Add(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			ErrNo:  500,
			ErrMsg: "添加失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   label,
		ErrNo:  0,
		ErrMsg: "",
	})
}

// Edit 编辑标签
func (h *LabelHandler) Edit(c *gin.Context) {
	var dto domain.LabelEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			ErrNo:  400,
			ErrMsg: "参数验证失败: " + err.Error(),
		})
		return
	}

	label, err := h.labelService.Edit(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			ErrNo:  500,
			ErrMsg: "编辑失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   label,
		ErrNo:  0,
		ErrMsg: "",
	})
}

// Delete 删除标签
func (h *LabelHandler) Delete(c *gin.Context) {
	var dto domain.LabelDeleteDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, domain.StandardResponse{
			ErrNo:  400,
			ErrMsg: "参数验证失败: " + err.Error(),
		})
		return
	}

	err := h.labelService.Delete(c.Request.Context(), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.StandardResponse{
			ErrNo:  500,
			ErrMsg: "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.StandardResponse{
		Data:   nil,
		ErrNo:  0,
		ErrMsg: "",
	})
}
