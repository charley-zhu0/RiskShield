package domain

// ModerationRequest 审核请求
type ModerationRequest struct {
	Model    string `form:"model" binding:"required"`
	Content  string `form:"content" binding:"required,max=50000"`
	TraceID  string `form:"trace_id"`
	App      string `form:"app"`
	Location string `form:"location"`
}

// ModerationResult 审核结果
type ModerationResult struct {
	Decision   string   `json:"decision"`    // PASS 或 REJECT
	Reason     string   `json:"reason"`      // 拒绝原因
	HitWords   []string `json:"hit_words"`   // 命中的敏感词
	SafetyType string   `json:"safety_type"` // 安全分类
	AttackType string   `json:"attack_type"` // 攻击类型
}
