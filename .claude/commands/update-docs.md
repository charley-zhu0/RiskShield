# Go 文档更新

将文档与 Go 代码库同步，从唯一真实来源文件生成文档。

## 步骤 1：识别唯一真实来源

| 来源 | 生成内容 |
|------|----------|
| `Makefile` / `Taskfile.yml` | 可用命令参考 |
| `config` 结构体 / `.env.example` | 环境变量文档 |
| `swagger.yaml` / `docs/swagger` | API 端点参考 |
| Godoc 注释（公开导出） | 公共 API 文档 |
| `go.mod` | 依赖和版本信息 |
| `Dockerfile` / `docker-compose.yml` | 基础设施配置文档 |

## 步骤 2：生成 Make/Task 命令参考

1. 读取 `Makefile`（或 `Taskfile.yml`）
2. 提取 target 及其注释/描述
3. 生成参考表格：

```markdown
| 命令 | 描述 |
|------|------|
| `make run` | 本地运行应用程序 |
| `make build` | 构建生产环境二进制文件 |
| `make test` | 运行所有测试并生成覆盖率报告 |
| `make lint` | 运行 golangci-lint |
```

## 步骤 3：生成配置文档

1. 查找配置结构体（通常在 `internal/config` 或 `cmd/...` 中）
2. 查找结构体标签，如 `env:"..."`、`default:"..."`、`json:"..."`
3. 如果存在 `.env.example` 也需检查
4. 记录预期格式、必填字段和默认值

```markdown
| 环境变量 | 配置结构体字段 | 必填 | 默认值 | 描述 |
|----------|----------------|------|--------|------|
| `DATABASE_URL` | `Database.URL` | 是 | - | PostgreSQL 连接字符串 |
| `HTTP_PORT` | `Server.Port` | 否 | 8080 | 监听端口 |
| `LOG_LEVEL` | `Log.Level` | 否 | info | 日志级别 |
```

## 步骤 4：更新贡献指南

生成或更新 `docs/CONTRIBUTING.md`，包含：
- Go 版本要求（从 `go.mod` 提取）
- 设置步骤（`go mod download` 等）
- 可用的 Make/Task 命令
- 测试流程（`go test ./...`）
- 代码风格（引用 `golangci-lint` 配置）
- 项目布局说明（介绍 `internal`、`pkg`、`cmd` 目录）

## 步骤 5：更新运行手册

生成或更新 `docs/RUNBOOK.md`，包含：
- 构建说明（`go build`）
- 二进制执行标志和参数
- 健康检查端点
- 常见运行时错误（panic、OOM）
- 可观测性（metrics 端口、pprof 端点）

## 步骤 6：API 文档检查

1. 检查是否配置了 `swag` 或 `openapi` 生成
2. 如果使用 `swag`，运行 `swag init` 更新 `docs/` 目录
3. 验证 `README.md` 是否链接到最新的 API 文档

## 步骤 7：过期检查

1. 查找 90 天以上未修改的文档文件
2. 与相关包中最近的 `.go` 文件变更进行交叉对比
3. 标记可能过时的文档

## 步骤 8：显示摘要

```
文档更新 (Go)
──────────────────────────────
已更新:  docs/CONTRIBUTING.md (Make targets)
已更新:  docs/CONFIG.md (在 internal/config 中检测到新环境变量)
已标记:  docs/DEPLOY.md (142 天未更新)
已跳过:  docs/API.md (swag init 无变更)
──────────────────────────────
```

## 规则

- **唯一真实来源**：优先从代码（结构体标签、Makefile）提取，而非手动文件。
- **保留手动内容**：仅更新生成的部分；保持手写内容不变。
- **标记生成内容**：在生成内容周围使用 `<!-- AUTO-GENERATED -->` 标记。
- **遵循 Go 惯例**：确保文档体现惯用的 Go 模式（如标准项目布局）。