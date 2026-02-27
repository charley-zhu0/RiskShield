package domain

import (
	"time"

	"github.com/google/uuid"
)

// Label 标签领域模型
type Label struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"column:tenantID;type:varchar(50);not null" json:"tenant_id"`
	Class     string    `gorm:"type:varchar(50);not null" json:"class"`
	SID       string    `gorm:"column:sid;type:varchar(50);not null;uniqueIndex" json:"sid"`
	PID       string    `gorm:"column:pid;type:varchar(50);not null" json:"pid"`
	Title     string    `gorm:"type:varchar(250);not null" json:"title"`
	Value     string    `gorm:"type:varchar(250);not null" json:"value"`
	Step      int       `gorm:"type:tinyint;not null" json:"step"`
	Source    int       `gorm:"type:tinyint;not null" json:"source"`
	CreateBy  string    `gorm:"column:createby;type:varchar(50)" json:"create_by"`
	UpdateBy  string    `gorm:"column:updateby;type:varchar(50)" json:"update_by"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Label) TableName() string {
	return "relate_list"
}

// GenerateSID 生成标签 SID (格式: LB-{uuid})
func GenerateSID() string {
	return "LB-" + uuid.New().String()
}

// CalculateStep 计算标签层级
// 如果 pid == "zhinaolabel"，则 step = 1
// 否则查询父标签，step = parent.Step + 1
func CalculateStep(pid string, parentStep int) int {
	if pid == "zhinaolabel" {
		return 1
	}
	return parentStep + 1
}

// LabelQueryDTO 查询请求 DTO
type LabelQueryDTO struct {
	Page     int    `form:"page" binding:"required,min=1"`
	Size     int    `form:"size" binding:"required,min=1,max=100"`
	SelfOnly *bool  `form:"self_only"`
	Locale   string `form:"locale" binding:"max=50"`
	PID      string `form:"pid" binding:"max=50"`
	Title    string `form:"title" binding:"max=250"`
	SID      string `form:"sid" binding:"max=50"`
}

// LabelAddDTO 添加请求 DTO
type LabelAddDTO struct {
	PID   string `json:"pid" binding:"required,max=50"`
	Title string `json:"title" binding:"required,min=1,max=250"`
}

// LabelEditDTO 编辑请求 DTO
type LabelEditDTO struct {
	ID    uint   `json:"id" binding:"required,min=1"`
	PID   string `json:"pid" binding:"required,max=50"`
	Title string `json:"title" binding:"required,min=1,max=250"`
}

// LabelDeleteDTO 删除请求 DTO
type LabelDeleteDTO struct {
	ID uint `json:"id" binding:"required"`
}

// LabelQueryResponse 查询响应
type LabelQueryResponse struct {
	Data  []Label `json:"data"`
	Total int64   `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
}
