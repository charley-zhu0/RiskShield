---
description: 分析 Go 项目的测试覆盖率，识别差距，并生成缺失的测试以达到 80% 以上的覆盖率。
---

# Go 测试覆盖率

分析 Go 项目的测试覆盖率，识别差距，并生成缺失的测试以达到 80% 以上的覆盖率。

## 步骤 1：运行覆盖率分析

### 基本覆盖率命令

```bash
# 基本覆盖率报告
go test -cover ./...

# 生成覆盖率 profile 文件
go test -coverprofile=coverage.out ./...

# 带竞态检测的覆盖率
go test -race -coverprofile=coverage.out ./...

# 原子模式覆盖率（推荐用于并发代码）
go test -covermode=atomic -coverprofile=coverage.out ./...
```

### 覆盖率模式

| 模式 | 说明 | 使用场景 |
|------|------|----------|
| `set` | 语句是否被执行（布尔值） | 默认模式，最快 |
| `count` | 语句执行次数 | 热点代码分析 |
| `atomic` | 并发安全的计数 | 并发代码测试 |

## 步骤 2：查看覆盖率报告

### 终端查看

```bash
# 按函数查看覆盖率
go tool cover -func=coverage.out

# 输出示例
# github.com/user/project/pkg/auth/auth.go:25:     ValidateToken           100.0%
# github.com/user/project/pkg/auth/auth.go:45:     CreateToken             66.7%
# github.com/user/project/pkg/user/user.go:12:     CreateUser              45.5%
# total:                                           (statements)            72.3%
```

### 浏览器查看（可视化）

```bash
# 生成 HTML 报告并在浏览器中打开
go tool cover -html=coverage.out

# 或保存为 HTML 文件
go tool cover -html=coverage.out -o coverage.html
```

### JSON 格式输出（CI/CD 集成）

```bash
# 使用第三方工具输出 JSON
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest

gocov convert coverage.out | gocov-xml > coverage.xml
```

## 步骤 3：分析覆盖率差距

### 识别低覆盖率文件

1. 运行覆盖率命令
2. 查看 `go tool cover -func` 输出
3. 列出**覆盖率低于 80%** 的文件，按最差的排序
4. 对于每个覆盖率不足的文件，识别：
   - 未测试的函数或方法
   - 缺失的分支覆盖（if/else、switch、错误路径）
   - 增加分母的死代码

### 常见的覆盖率盲区

```go
// 1. 错误处理分支
if err != nil {
    return err  // ← 常常未被测试
}

// 2. switch 语句的 default 分支
switch v {
case "a":
    // tested
case "b":
    // tested
default:
    // ← 常常未被测试
}

// 3. 提前返回的分支
if !valid {
    return nil, ErrInvalid  // ← 需要测试无效输入
}

// 4. 并发代码的错误路径
select {
case result := <-ch:
    // tested
case <-ctx.Done():
    // ← 超时路径常常未被测试
}
```

## 步骤 4：生成缺失的测试

对于每个覆盖率不足的文件，按以下优先级生成测试：

### 测试优先级

1. **正常路径** — 使用有效输入的核心功能
2. **错误处理** — 无效输入、缺失数据、网络失败
3. **边缘情况** — 空切片、nil、边界值（0、-1、MAX_INT）
4. **分支覆盖** — 每个 if/else、switch case

### 表格驱动测试模板

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr error
    }{
        // 正常路径
        {
            name:  "valid input returns expected result",
            input: validInput,
            want:  expectedOutput,
        },
        // 错误处理
        {
            name:    "empty input returns error",
            input:   "",
            wantErr: ErrEmptyInput,
        },
        // 边缘情况
        {
            name:  "nil slice returns empty result",
            input: nil,
            want:  emptyResult,
        },
        {
            name:  "max value boundary",
            input: math.MaxInt64,
            want:  boundaryResult,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            
            if tt.wantErr != nil {
                if !errors.Is(err, tt.wantErr) {
                    t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            }
            
            if err != nil {
                t.Errorf("Function() unexpected error = %v", err)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 测试生成规则

- 测试文件放在相同目录：`foo.go` → `foo_test.go`
- 遵循项目现有的测试模式
- 模拟外部依赖（数据库、API、文件系统）
- 每个测试应独立 — 测试之间无共享可变状态
- 描述性命名：`TestValidateEmail_EmptyInput_ReturnsError`

## 步骤 5：验证改进

```bash
# 运行完整测试套件
go test ./...

# 重新运行覆盖率分析
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# 如果仍低于 80%，重复步骤 4
```

## 步骤 6：报告

显示前后对比：

```
Go 覆盖率报告
══════════════════════════════════════════════════════
文件                                    之前    之后
pkg/auth/token.go                       45%     88%
pkg/user/validation.go                  32%     82%
internal/service/order.go               67%     91%
══════════════════════════════════════════════════════
整体覆盖率:                              67%     84%  ✅
```

## 覆盖率目标

| 代码类型 | 目标覆盖率 |
|---------|-----------|
| 关键业务逻辑 | 100% |
| 公共 API（导出函数） | 90%+ |
| 通用代码 | 80%+ |
| 生成代码（protobuf 等） | 排除 |

## 排除生成代码

### 使用构建标签排除

```go
//go:build !generate
// +build !generate

package mypackage

// 生成的代码...
```

### 在 CI 中排除特定文件

```bash
# 排除 mock 和生成的文件
go test -coverprofile=coverage.out ./... \
    -coverpkg=./... \
    | grep -v "_mock.go" \
    | grep -v ".pb.go"
```

## 重点关注区域

- 高圈复杂度的函数（复杂分支逻辑）
- 错误处理和 recover 块
- 跨代码库使用的工具函数
- HTTP/gRPC 处理器（请求 → 响应流程）
- 并发代码（goroutine、channel、mutex）
- 边缘情况：nil、空字符串、空切片、零、负数

## 常用工具

```bash
# gocov - 覆盖率分析工具
go install github.com/axw/gocov/gocov@latest

# gocov-xml - XML 格式输出（CI 集成）
go install github.com/AlekSi/gocov-xml@latest

# gocov-html - HTML 报告
go install github.com/matm/gocov-html/cmd/gocov-html@latest

# gotestsum - 更好的测试输出格式
go install gotest.tools/gotestsum@latest

# 示例：生成详细报告
gotestsum --format testname -- -coverprofile=coverage.out ./...
gocov convert coverage.out | gocov-html > coverage.html
```

## CI/CD 集成示例

### GitHub Actions

```yaml
name: Test Coverage

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run tests with coverage
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
      
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage is below 80%"
            exit 1
          fi
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          files: coverage.out
```

## 相关命令

- `/tdd` - Go TDD 工作流
- `/go-build-fix` - 修复构建错误
- `/go-review` - 代码质量审查
- `/verify` - 完整验证循环

## 相关资源

- 技能: `skills/golang-testing/`
- 技能: `skills/golang-patterns/`