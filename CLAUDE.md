# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

RiskShield 是基于 Go + Gin 的文本内容审核服务，通过敏感词检测和大模型分类实现风险评估。

## 核心架构

采用分层架构：
- **cmd/server**: 应用入口，依赖注入和路由配置
- **internal/domain**: 领域模型（ModerationRequest, ModerationResult）
- **internal/client**: 外部服务客户端（敏感词服务、LLM服务）
- **internal/service**: 业务逻辑层，实现审核决策逻辑
- **internal/handler**: HTTP Handler，处理 Gin 请求
- **internal/config**: 环境变量配置管理

审核决策流程：
1. 敏感词检测命中 → 直接 REJECT
2. 大模型分类 safety_type != "正常" → REJECT
3. 两者都通过 → PASS

## 开发命令

```bash
# 安装依赖
go mod tidy

# 运行所有测试
go test ./...

# 测试覆盖率（要求 ≥80%）
go test -cover ./internal/...
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out

# 运行单个测试
go test -v -run TestFunctionName ./internal/package

# 并发安全检测
go test -race ./...

# 编译
go build -o bin/server ./cmd/server

# 运行服务
export SERVER_PORT=8080
export SENSITIVE_WORD_URL=http://localhost:8081
export LLM_SERVICE_URL=http://localhost:8082
./bin/server
```

## 开发规范

- 严格遵循 TDD：先写表驱动测试，再实现功能
- 测试覆盖率必须 ≥80%
- 使用 Mock 隔离外部依赖
- 所有测试必须通过 race 检测
- 接口定义在 service 和 client 层，便于测试注入

## 环境变量

- `SERVER_PORT`: 服务端口（默认 8080）
- `SENSITIVE_WORD_URL`: 敏感词服务地址
- `LLM_SERVICE_URL`: 大模型服务地址
- `REQUEST_TIMEOUT`: 请求超时时间（秒，默认 30）
