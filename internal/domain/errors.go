package domain

import "errors"

var (
	ErrInvalidRequest    = errors.New("无效的请求参数")
	ErrSensitiveWord     = errors.New("敏感词服务调用失败")
	ErrLLMService        = errors.New("大模型服务调用失败")
	ErrTimeout           = errors.New("请求超时")
)
