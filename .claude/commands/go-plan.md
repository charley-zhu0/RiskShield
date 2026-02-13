---
description: 为 Go 项目重新阐述需求、评估风险并创建分步实施计划。在接触任何代码之前，等待用户确认。
---

# Go 计划命令

此命令调用 **planner** 代理，在编写任何 Go 代码之前创建一个全面的实施计划。

## 此命令的作用

1. **重新阐述需求** - 明确需要构建什么（Go 项目特定）
2. **识别风险** - 揭示潜在问题和阻碍（Go 并发、错误处理等）
3. **创建分步计划** - 将实施分解为多个阶段（遵循 Go 最佳实践）
4. **等待确认** - 必须获得用户批准才能继续

## 何时使用

在以下情况下使用 `/go-plan`：

* 开始新的 Go 功能时
* 进行重大 Go 架构变更时（例如：重构包结构、引入新的依赖）
* 处理复杂的 Go 重构时（例如：改善错误处理、优化并发模型）
* 多个 Go 文件/包将受到影响时
* 需要明确 Go 项目的依赖关系和接口设计时

## 工作原理

planner 代理将：

1. **分析请求** 并用清晰的术语重新阐述需求（Go 特定术语）
2. **分解为多个阶段**，每个阶段包含具体、可操作的步骤（Go 文件路径）
3. **识别组件之间的依赖关系**（Go 包依赖、接口实现）
4. **评估风险和潜在阻碍**（并发问题、错误处理、性能瓶颈）
5. **估算复杂度**（高/中/低）
6. **展示计划**并等待您明确的确认

## 使用示例

```
User: /go-plan 我需要实现一个 HTTP 服务，支持用户注册和登录功能，使用 PostgreSQL 存储用户数据

Agent (planner):
# 实施计划: 用户认证服务

## 概述
实现用户注册和登录功能，包括密码哈希存储、JWT 令牌生成和验证。使用 PostgreSQL 作为数据存储。

## 需求
- 用户可以通过邮箱和密码注册
- 用户可以使用邮箱和密码登录
- 密码必须加盐哈希存储
- 登录成功后返回 JWT 令牌
- 邮箱必须唯一

## 架构变更
- 新增 `User` 结构体和 Repository 接口 (`internal/domain/user.go`)
- 新增 PostgreSQL 表 `users`
- 新增 HTTP Handler (`internal/handler/auth_handler.go`)
- 引入 `golang.org/x/crypto/bcrypt` 库用于密码哈希
- 引入 `github.com/golang-jwt/jwt/v5` 库用于 JWT 处理

## 实施步骤

### 阶段 1: 领域层与数据库 (Domain & DB)
1. **定义用户领域模型** (文件: internal/domain/user.go)
   - 行动: 定义 `User` struct, `RegisterRequest` struct, `LoginRequest` struct
   - 原因: 核心业务实体定义
   - 依赖: 无
   - 风险: 低

2. **定义 Repository 接口** (文件: internal/domain/user_repo.go)
   - 行动: 定义 `UserRepository` interface (Create, GetByEmail)
   - 原因: 依赖倒置，便于测试
   - 依赖: 步骤 1
   - 风险: 低

### 阶段 2: 数据库实现
3. **实现 PostgreSQL Repository** (文件: internal/repository/postgres/user_repo.go)
   - 行动: 实现 `UserRepository`，编写 SQL 查询，处理 `sql.ErrNoRows`
   - 原因: 数据持久化
   - 依赖: 步骤 2
   - 风险: 中 (需确保事务处理正确)

### 阶段 3: 服务层
4. **实现认证服务** (文件: internal/service/auth_service.go)
   - 行动: 实现注册和登录逻辑，调用 `UserRepository`，处理密码哈希和 JWT 生成
   - 原因: 封装业务规则
   - 依赖: 步骤 2
   - 风险: 高 (密码安全、JWT 安全)

### 阶段 4: 接口层
5. **实现 HTTP Handler** (文件: internal/handler/auth_handler.go)
   - 行动: 解析 JSON 请求，验证输入，调用 Service，返回响应
   - 原因: 暴露 API
   - 依赖: 步骤 4
   - 风险: 中 (输入验证)

6. **注册路由** (文件: cmd/server/main.go)
   - 行动: 将 `/register` 和 `/login` 路由绑定到 Handler
   - 原因: 启用端点
   - 依赖: 步骤 5

## 测试策略
- 单元测试: `auth_service_test.go` (使用表格驱动测试，Mock Repository)
- 集成测试: `user_repo_test.go` (使用测试数据库)
- E2E 测试: HTTP 请求测试注册和登录流程

## 风险与缓解
- **风险**: 密码哈希过慢影响 API 性能
  - 缓解: 调整哈希参数成本，考虑异步处理非关键路径
- **风险**: JWT 密钥泄露
  - 缓解: 使用环境变量存储密钥，定期轮换

## 成功标准
- [ ] 能够成功注册新用户
- [ ] 能够成功登录并获取 JWT
- [ ] 数据库中密码已加密
- [ ] 重复邮箱注册返回 409 错误
- [ ] 单元测试覆盖率 > 80%
- [ ] `golangci-lint` 检查通过

**WAITING FOR CONFIRMATION**: Proceed with this plan? (yes/no/modify)
```

## 重要说明

**关键**：planner 代理在您明确用"是"、"继续"或类似的肯定性答复确认计划之前，**不会**编写任何代码。

如果您希望修改，请回复：

* "修改：\[您的修改内容]"
* "不同方法：\[替代方案]"
* "跳过阶段 2，先执行阶段 3"

## 与其他命令的集成

计划之后：

* 使用 `/tdd` 以测试驱动开发的方式实施
* 如果出现构建错误，使用 `/build-fix`
* 使用 `/go-review` 审查已完成的实施

## 相关代理

此命令调用位于 `.claude/agents/planner.md` 的 `planner` 代理。