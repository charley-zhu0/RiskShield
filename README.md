# RiskShield - 文本检测服务

基于 Go + Gin 框架的文本内容审核服务，通过敏感词检测和大模型分类实现综合风险评估。

## 项目结构

```
RiskShield/
├── cmd/server/main.go              # 应用入口
├── internal/
│   ├── config/                     # 配置管理
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                     # 领域模型
│   │   ├── moderation.go
│   │   └── errors.go
│   ├── client/                     # 外部服务客户端
│   │   ├── sensitiveword.go
│   │   ├── sensitiveword_test.go
│   │   ├── llm.go
│   │   └── llm_test.go
│   ├── service/                    # 业务逻辑层
│   │   ├── moderation.go
│   │   └── moderation_test.go
│   └── handler/                    # HTTP Handler
│       ├── moderation.go
│       └── moderation_test.go
├── go.mod
└── go.sum
```

## 功能特性

- ✅ 敏感词检测（调用敏感词服务）
- ✅ 大模型内容分类（调用 LLM 服务）
- ✅ 综合风险决策（PASS/REJECT）
- ✅ 完整的单元测试和集成测试
- ✅ 87.7% 测试覆盖率
- ✅ 并发安全（通过 race 检测）

## 环境变量

```bash
SERVER_PORT=8080                                    # 服务端口
SENSITIVE_WORD_URL=http://localhost:8081           # 敏感词服务地址
LLM_SERVICE_URL=http://localhost:8082              # 大模型服务地址
REQUEST_TIMEOUT=30                                  # 请求超时时间（秒）
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./internal/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

### 3. 编译

```bash
go build -o bin/server ./cmd/server
```

### 4. 运行

```bash
export SERVER_PORT=8080
export SENSITIVE_WORD_URL=http://localhost:8081
export LLM_SERVICE_URL=http://localhost:8082
./bin/server
```

## API 接口

### POST /v1/private/moderations/text

审核文本内容

**请求参数** (multipart/form-data):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| model | string | 是 | 模型名称 |
| content | string | 是 | 待审核内容 |
| trace_id | string | 否 | 追踪ID |
| app | string | 否 | 应用标识 |
| location | string | 否 | 位置信息 |

**响应示例**:

```json
{
  "decision": "REJECT",
  "reason": "命中敏感词",
  "hit_words": ["敏感词"],
  "safety_type": "涉政:敏感政治事件",
  "attack_type": "无"
}
```

## 测试覆盖率

| 模块 | 覆盖率 |
|------|--------|
| client | 83.3% |
| config | 100.0% |
| handler | 80.0% |
| service | 100.0% |
| **总计** | **87.7%** |

## 决策逻辑

1. **敏感词检测**: 命中敏感词 → 直接 REJECT
2. **大模型分类**: safety_type != "正常" → REJECT
3. **两者都通过**: → PASS

## 技术栈

- Go 1.21+
- Gin Web Framework v1.10.0
- 标准库 testing 包
- httptest 用于 HTTP 测试

## 开发规范

- ✅ 严格遵循 TDD 开发流程
- ✅ 使用表驱动测试
- ✅ Mock 外部依赖
- ✅ 测试覆盖率 ≥ 80%
- ✅ 通过 race 检测
