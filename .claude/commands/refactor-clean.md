---
description: Safely identify and remove dead Go code with test verification at every step. Invokes static analysis tools and enforces safe deletion practices.
---

# Go Refactor Clean

安全识别并移除 Go 项目中的死代码，每一步都通过测试验证。

## 步骤 1：检测死代码

运行 Go 专用的静态分析工具：

| 工具 | 检测内容 | 命令 |
|------|----------|------|
| deadcode | 未使用的函数和方法 | `go install golang.org/x/tools/cmd/deadcode@latest && deadcode ./...` |
| staticcheck | 未使用的代码、废弃 API | `staticcheck ./...` |
| golangci-lint | 综合检测（unused, deadcode, varcheck） | `golangci-lint run --enable=unused,deadcode` |
| go vet | 基础静态分析 | `go vet ./...` |
| unparam | 未使用的函数参数 | `go install mvdan.cc/unparam@latest && unparam ./...` |

### 推荐的检测流程

```bash
# 1. 安装工具（如果尚未安装）
go install golang.org/x/tools/cmd/deadcode@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# 2. 运行死代码检测
deadcode -test ./...

# 3. 运行综合静态分析
staticcheck ./...

# 4. 如果有 golangci-lint，运行完整检查
golangci-lint run
```

### 手动检测方法

如果工具不可用，使用 grep 查找可能的死代码：

```bash
# 查找未导出的函数，检查是否被引用
grep -rn "^func [a-z]" --include="*.go" .

# 查找定义但可能未使用的类型
grep -rn "^type [A-Z]" --include="*.go" .
```

## 步骤 2：分类发现的问题

按安全级别对发现的问题进行分类：

| 级别 | 示例 | 操作 |
|------|------|------|
| **安全** | 未使用的内部函数、测试辅助函数、私有方法 | 可以放心删除 |
| **谨慎** | 未使用的导出函数、HTTP 处理器、中间件 | 验证无反射调用或外部消费者 |
| **危险** | init() 函数、接口定义、配置结构体 | 删除前需深入调查 |

### Go 特有的注意事项

- **init() 函数**：即使看起来未被调用，也会在包初始化时自动执行
- **接口实现**：检查是否有类型隐式实现了该接口
- **反射调用**：搜索 `reflect.` 使用，可能动态调用函数
- **CGO 导出**：带有 `//export` 注释的函数可能被 C 代码调用
- **插件系统**：检查是否使用 `plugin` 包动态加载

## 步骤 3：安全删除循环

对于每个 **安全** 级别的项目：

1. **运行完整测试套件** — 建立基线（全部通过）
   ```bash
   go test ./... -race -count=1
   ```

2. **删除死代码** — 使用编辑工具进行精确删除

3. **重新运行测试套件** — 验证没有破坏任何功能
   ```bash
   go test ./... -race -count=1
   ```

4. **如果测试失败** — 立即回滚并跳过此项
   ```bash
   git checkout -- <file>
   ```

5. **如果测试通过** — 继续下一项

6. **运行构建验证** — 确保编译通过
   ```bash
   go build ./...
   ```

## 步骤 4：处理谨慎级别的项目

删除谨慎级别的项目前，执行以下检查：

```bash
# 搜索反射调用
grep -rn "reflect\." --include="*.go" .

# 搜索字符串形式的函数名引用（可能用于动态调用）
grep -rn "\"FunctionName\"" --include="*.go" .

# 检查是否在 HTTP 路由中注册
grep -rn "HandleFunc\|Handle\|Router" --include="*.go" .

# 检查是否在 gRPC 服务中注册
grep -rn "RegisterServer\|pb\." --include="*.go" .

# 检查是否有外部包依赖此导出
go list -m all | head -20  # 查看依赖关系
```

### 导出函数的额外检查

对于导出函数（首字母大写），需要确认：

- 不是公共 API 的一部分
- 不被其他模块导入使用
- 不通过反射调用
- 不是接口实现的一部分

## 步骤 5：合并重复代码

删除死代码后，查找以下问题：

| 问题类型 | 检测方法 | 操作 |
|----------|----------|------|
| 近似重复函数（>80% 相似） | 代码审查 | 合并为一个通用函数 |
| 冗余类型定义 | `grep -rn "^type"` | 合并或使用类型别名 |
| 无价值的包装函数 | 代码审查 | 内联调用 |
| 无意义的重导出 | 检查 `type X = Y` | 移除间接层 |

### Go 特有的重构建议

```go
// 不好：无价值的包装
func GetUserByID(id string) (*User, error) {
    return db.FindUser(id)
}

// 好：直接使用 db.FindUser，删除包装函数

// 不好：重复的错误处理
if err != nil {
    return fmt.Errorf("failed to get user: %w", err)
}
// ... 多处相同模式

// 好：考虑使用辅助函数或统一错误处理
```

## 步骤 6：验证与报告

### 最终验证

```bash
# 完整构建
go build ./...

# 完整测试（含竞态检测）
go test ./... -race -cover

# 静态分析
go vet ./...
staticcheck ./...

# 检查是否引入新问题
golangci-lint run
```

### 报告结果

```
Go 死代码清理报告
══════════════════════════════════
已删除:   8 个未使用的函数
          2 个未使用的类型
          3 个未使用的常量
          1 个未使用的文件
已跳过:   2 项（测试失败）
          1 项（可能有反射调用）
节省:     约 320 行代码
══════════════════════════════════
所有测试通过 ✅
构建成功 ✅
静态分析通过 ✅
```

## 规则

- **删除前必须运行测试** — 永远不要在没有测试验证的情况下删除代码
- **一次删除一项** — 原子性更改使回滚更容易
- **不确定时跳过** — 保留死代码比破坏生产环境更好
- **清理时不要重构** — 分离关注点（先清理，后重构）
- **注意 init() 函数** — Go 的 init() 会自动执行，即使看起来未被调用
- **检查接口实现** — 删除类型前确认没有隐式接口实现
- **保留必要的空白导入** — `import _ "package"` 可能用于副作用

## 常用工具安装

```bash
# deadcode - 官方死代码检测工具
go install golang.org/x/tools/cmd/deadcode@latest

# staticcheck - 综合静态分析
go install honnef.co/go/tools/cmd/staticcheck@latest

# golangci-lint - 多 linter 聚合
# macOS
brew install golangci-lint
# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# unparam - 未使用参数检测
go install mvdan.cc/unparam@latest
```

## 相关命令

- `/tdd` - 运行测试验证
- `/go-build-fix` - 修复构建错误
- `/go-review` - 代码质量审查
- `/verify` - 完整验证循环

## 相关资源

- Agent: `agents/go-reviewer.md`
- Skill: `skills/golang-patterns/`
- Skill: `skills/golang-testing/`