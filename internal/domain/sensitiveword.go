package domain

import "time"

// SensitiveWord 敏感词领域模型
type SensitiveWord struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    string    `gorm:"column:tenantID;type:varchar(50);not null" json:"tenantID"`
	App         string    `gorm:"type:varchar(50)" json:"app"`
	Location    string    `gorm:"type:varchar(50)" json:"location"`
	Word        string    `gorm:"type:varchar(255);not null" json:"word"`
	FirstLabel  string    `gorm:"column:firstLabel;type:varchar(50)" json:"firstLabel"`
	SecondLabel string    `gorm:"column:secondLabel;type:varchar(50)" json:"secondLabel"`
	ThirdLabel  string    `gorm:"column:thirdLabel;type:varchar(50)" json:"thirdLabel"`
	MatchedID   int       `gorm:"column:matchedId;type:int" json:"matchedId"`
	QueryDeal   int       `gorm:"column:queryDeal;type:int" json:"queryDeal"`
	Source      int       `gorm:"type:tinyint" json:"source"`
	CreateBy    string    `gorm:"column:createby;type:varchar(50)" json:"createby"`
	UpdateBy    string    `gorm:"column:updateby;type:varchar(50)" json:"updateby"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (SensitiveWord) TableName() string {
	return "sensitiveword"
}

// SensitiveWordQueryDTO 查询请求 DTO
type SensitiveWordQueryDTO struct {
	Page        int    `form:"page" binding:"required,min=1"`
	Size        int    `form:"size" binding:"required,min=1,max=100"`
	Word        string `form:"word" binding:"max=255"`
	FirstLabel  string `form:"firstLabel" binding:"max=50"`
	SecondLabel string `form:"secondLabel" binding:"max=50"`
	QueryDeal   *[]int `form:"queryDeal"`
	SelfOnly    bool   `form:"self_only"`
	Locale      string `form:"locale" binding:"max=10"`
}

// SensitiveWordAddDTO 批量添加请求 DTO
type SensitiveWordAddDTO struct {
	Words       []string `json:"words" binding:"required,min=1,dive,required,min=1,max=255"`
	FirstLabel  string   `json:"firstLabel" binding:"required,max=50"`
	SecondLabel string   `json:"secondLabel" binding:"required,max=50"`
	ThirdLabel  string   `json:"thirdLabel" binding:"max=50"`
	QueryDeal   int      `json:"queryDeal" binding:"required,min=0"`
	MatchedID   int      `json:"matchedId" binding:"min=0"`
}

// SensitiveWordEditDTO 编辑请求 DTO
type SensitiveWordEditDTO struct {
	ID          uint   `json:"id" binding:"required,min=1"`
	Word        string `json:"word" binding:"required,min=1,max=255"`
	FirstLabel  string `json:"firstLabel" binding:"required,max=50"`
	SecondLabel string `json:"secondLabel" binding:"required,max=50"`
	ThirdLabel  string `json:"thirdLabel" binding:"max=50"`
	QueryDeal   int    `json:"queryDeal" binding:"required,min=0"`
	MatchedID   int    `json:"matchedId" binding:"min=0"`
}

// SensitiveWordDeleteDTO 删除请求 DTO
type SensitiveWordDeleteDTO struct {
	ID uint `json:"id" binding:"required,min=1"`
}

// SensitiveWordQueryResponse 查询响应
type SensitiveWordQueryResponse struct {
	List  []SensitiveWord `json:"list"`
	Total int64           `json:"total"`
}
