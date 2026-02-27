package domain

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

// Intervention 干预库领域模型
type Intervention struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"column:tenantID;type:varchar(50);not null" json:"tenant_id"`
	Query     string    `gorm:"type:text;not null" json:"query"`
	Answer    string    `gorm:"type:text;not null" json:"answer"`
	QueryHash string    `gorm:"column:queryHash;type:varchar(50);not null" json:"query_hash"`
	Source    int       `gorm:"type:tinyint;not null" json:"source"`
	CreateBy  string    `gorm:"column:createby;type:varchar(50)" json:"create_by"`
	UpdateBy  string    `gorm:"column:updateby;type:varchar(50)" json:"update_by"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Intervention) TableName() string {
	return "intervention"
}

// CalculateQueryHash 计算 query 的 MD5 哈希值
func CalculateQueryHash(query string) string {
	hash := md5.Sum([]byte(query))
	return hex.EncodeToString(hash[:])
}

// InterventionQueryDTO 查询请求 DTO
type InterventionQueryDTO struct {
	Page   int    `form:"page" binding:"required,min=1"`
	Size   int    `form:"size" binding:"required,min=1,max=100"`
	Query  string `form:"query" binding:"max=500"`
	Source *int   `form:"source" binding:"omitempty,min=0"`
	Fuzzy  bool   `form:"fuzzy"`
}

// InterventionAddDTO 添加请求 DTO
type InterventionAddDTO struct {
	Query  string `json:"query" binding:"required,min=1,max=2000"`
	Answer string `json:"answer" binding:"required,min=1,max=5000"`
	Source int    `json:"source" binding:"required,min=0,max=10"`
}

// InterventionEditDTO 编辑请求 DTO
type InterventionEditDTO struct {
	ID     uint   `json:"id" binding:"required,min=1"`
	Query  string `json:"query" binding:"required,min=1,max=2000"`
	Answer string `json:"answer" binding:"required,min=1,max=5000"`
	Source int    `json:"source" binding:"required,min=0,max=10"`
}

// InterventionDeleteDTO 删除请求 DTO
type InterventionDeleteDTO struct {
	ID uint `json:"id" binding:"required"`
}

// InterventionQueryResponse 查询响应
type InterventionQueryResponse struct {
	Data  []Intervention `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// StandardResponse 标准响应格式
type StandardResponse struct {
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errmsg"`
	ErrNo  int         `json:"errno"`
}
