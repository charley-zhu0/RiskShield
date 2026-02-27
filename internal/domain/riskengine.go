package domain

import "time"

// RiskEngine 防护策略领域模型
type RiskEngine struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    string    `gorm:"column:tenantID;type:varchar(50);not null" json:"tenant_id"`
	App         string    `gorm:"type:varchar(100);not null" json:"app"`
	Location    string    `gorm:"type:varchar(100);not null" json:"location"`
	FirstLabel  string    `gorm:"column:firstLabel;type:varchar(100);not null" json:"first_label"`
	SecondLabel string    `gorm:"column:secondLabel;type:varchar(100);not null" json:"second_label"`
	ThirdLabel  string    `gorm:"column:thirdLabel;type:varchar(100);not null" json:"third_label"`
	QueryDeal   int       `gorm:"column:queryDeal;type:tinyint;not null" json:"query_deal"` // 0放行/1拦截/2审核
	Source      int       `gorm:"type:tinyint;not null" json:"source"`                      // 1系统/2用户
	CreateBy    string    `gorm:"column:createby;type:varchar(50)" json:"create_by"`
	UpdateBy    string    `gorm:"column:updateby;type:varchar(50)" json:"update_by"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (RiskEngine) TableName() string {
	return "riskengine"
}

// RiskEngineQueryDTO 查询请求 DTO
type RiskEngineQueryDTO struct {
	Page       int     `form:"page" binding:"required,min=1"`
	Size       int     `form:"size" binding:"required,min=1,max=100"`
	SelfOnly   bool    `form:"self_only"`
	Locale     *string `form:"locale" binding:"omitempty,max=100"`
	ThirdLabel *string `form:"thirdLabel" binding:"omitempty,max=100"`
}

// RiskEngineAddDTO 添加请求 DTO
type RiskEngineAddDTO struct {
	App         string `json:"app" binding:"required,min=1,max=100"`
	Location    string `json:"location" binding:"required,min=1,max=100"`
	FirstLabel  string `json:"first_label" binding:"required,min=1,max=100"`
	SecondLabel string `json:"second_label" binding:"required,min=1,max=100"`
	ThirdLabel  string `json:"third_label" binding:"required,min=1,max=100"`
	QueryDeal   int    `json:"query_deal" binding:"min=0,max=2"`
	Source      int    `json:"source" binding:"required,min=1,max=2"`
}

// RiskEngineEditDTO 编辑请求 DTO
type RiskEngineEditDTO struct {
	ID          uint   `json:"id" binding:"required,min=1"`
	App         string `json:"app" binding:"required,min=1,max=100"`
	Location    string `json:"location" binding:"required,min=1,max=100"`
	FirstLabel  string `json:"first_label" binding:"required,min=1,max=100"`
	SecondLabel string `json:"second_label" binding:"required,min=1,max=100"`
	ThirdLabel  string `json:"third_label" binding:"required,min=1,max=100"`
	QueryDeal   int    `json:"query_deal" binding:"min=0,max=2"`
	Source      int    `json:"source" binding:"required,min=1,max=2"`
}

// RiskEngineDeleteDTO 删除请求 DTO
type RiskEngineDeleteDTO struct {
	ID uint `json:"id" binding:"required"`
}

// RiskEngineQueryResponse 查询响应
type RiskEngineQueryResponse struct {
	Data  []RiskEngine `json:"data"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}
