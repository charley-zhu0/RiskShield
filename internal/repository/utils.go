package repository

import "strings"

// EscapeLikeString 转义 LIKE 查询中的特殊字符
// 防止 SQL 注入和意外的通配符匹配
func EscapeLikeString(pattern string) string {
	// 1. 先转义反斜杠
	pattern = strings.ReplaceAll(pattern, `\`, `\\`)
	// 2. 转义百分号
	pattern = strings.ReplaceAll(pattern, `%`, `\%`)
	// 3. 转义下划线
	pattern = strings.ReplaceAll(pattern, `_`, `\_`)
	return pattern
}
