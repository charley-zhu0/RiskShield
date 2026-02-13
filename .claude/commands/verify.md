# 验证命令 (Verification Command)

对当前 Go 代码库状态运行全面验证。

## 指令 (Instructions)

严格按照以下顺序执行验证：

1. **构建检查 (Build Check)**
   - 运行 `go build ./...` 构建所有包
   - 如果失败，报告错误并**停止**

2. **格式化检查 (Format Check)**
   - 运行 `goimports -l .` 检查代码格式
   - 报告所有未格式化的文件

3. **语法检查 (Vet Check)**
   - 运行 `go vet ./...` 检查代码问题
   - 报告所有发现的问题

4. **依赖检查 (Dependency Check)**
   - 运行 `go mod verify` 验证模块完整性
   - 运行 `go mod tidy -diff` 检查依赖是否整洁

5. **测试套件 (Test Suite)**
   - 运行 `go test ./... -cover -race` 执行所有测试（含覆盖率和竞态检测）
   - 报告通过/失败的数量
   - 报告覆盖率百分比（要求不低于 80%）

6. **调试语句审计 (Debug Statement Audit)**
   - 搜索源文件中的 `fmt.Println` 和 `print` 语句
   - 报告具体位置

7. **TODO/FIXME 审计 (TODO/FIXME Audit)**
   - 搜索源文件中的 TODO 和 FIXME 注释
   - 报告具体位置

8. **Git 状态 (Git Status)**
   - 显示未提交的更改
   - 显示自上次提交以来修改的文件

## 输出 (Output)

生成一份简明的验证报告：

```
验证结果: [通过/失败]

构建:    [正常/失败]
格式化:  [正常/X 个未格式化文件]
语法:    [正常/X 个问题]
依赖:    [正常/有问题]
测试:    [X/Y 通过, Z% 覆盖率]
调试:    [正常/发现 X 个调试语句]
待办:    [正常/发现 X 个 TODO/FIXME]

准备好提交 PR: [是/否]
```

如果有任何关键问题，请列出并附上修复建议。

## 参数 (Arguments)

$ARGUMENTS 可以是：
- `quick` - 仅运行构建 + 格式化检查
- `full` - 运行所有检查（默认）
- `pre-commit` - 运行提交前相关的检查
- `pre-pr` - 运行完整检查加上安全扫描
