---
name: tdd-workflow
description: 在编写新功能、修复错误或重构代码时使用此技能。强制执行测试驱动开发，要求 80% 以上的覆盖率，包括单元测试和集成测试。
---

# 测试驱动开发工作流

此技能确保所有代码开发都遵循 TDD 原则，并具有全面的测试覆盖率。

## 何时激活

- 编写新功能或功能实现
- 修复错误或问题
- 重构现有代码
- 添加 API 端点
- 创建新组件

## 核心原则

### 1. 测试优先
始终先编写测试，然后实现代码以使测试通过。

### 2. 覆盖率要求
- 最低 80% 覆盖率（单元测试 + 集成测试）
- 覆盖所有边界情况
- 测试错误场景
- 验证边界条件

### 3. 测试类型

#### 单元测试
- 单个函数和工具函数
- 组件逻辑
- 纯函数
- 辅助函数和工具函数

#### 集成测试
- API 端点
- 数据库操作
- 服务交互
- 外部 API 调用


## TDD 工作流步骤

### 步骤 1：编写用户旅程
```
作为一个 [角色]，我想要 [操作]，以便 [获得好处]

示例：
作为一个用户，我想要语义搜索市场，
以便即使没有精确的关键词也能找到相关的市场。
```

### 步骤 2：生成测试用例
为每个用户旅程创建全面的测试用例：

```go
func TestSemanticSearch(t *testing.T) {
    t.Run("返回查询的相关市场", func(t *testing.T) {
        // 测试实现
    })

    t.Run("优雅地处理空查询", func(t *testing.T) {
        // 测试边界情况
    })

    t.Run("当 Redis 不可用时回退到子字符串搜索", func(t *testing.T) {
        // 测试回退行为
    })

    t.Run("按相似度分数对结果排序", func(t *testing.T) {
        // 测试排序逻辑
    })
})
```

### 步骤 3：运行测试（它们应该失败）
```bash
go test ./...
# 测试应该失败 - 我们还没有实现
```

### 步骤 4：实现代码
编写最少的代码以使测试通过：

```go
// 由测试指导的实现
func SearchMarkets(query string) ([]Market, error) {
    // 实现代码
    return nil, nil
}
```

### 步骤 5：再次运行测试
```bash
go test ./...
# 测试现在应该通过
```

### 步骤 6：重构
在保持测试通过的同时提高代码质量：
- 删除重复代码
- 改进命名
- 优化性能
- 增强可读性

### 步骤 7：验证覆盖率
```bash
go test -cover ./...
# 验证达到 80% 以上的覆盖率
```

## 测试模式

### 单元测试模式（testing 包）
```go
package utils

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数相加", 2, 3, 5},
        {"负数相加", -2, -3, -5},
        {"零值相加", 0, 5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; 期望 %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

func TestValidateEmail(t *testing.T) {
    t.Run("有效邮箱", func(t *testing.T) {
        err := ValidateEmail("test@example.com")
        if err != nil {
            t.Errorf("期望 nil，得到 %v", err)
        }
    })

    t.Run("无效邮箱", func(t *testing.T) {
        err := ValidateEmail("invalid-email")
        if err == nil {
            t.Error("期望错误，得到 nil")
        }
    })
}
```

### API 集成测试模式
```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGetMarkets(t *testing.T) {
    t.Run("成功返回市场", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/markets", nil)
        w := httptest.NewRecorder()

        GetMarkets(w, req)

        if w.Code != http.StatusOK {
            t.Errorf("期望状态码 %d，得到 %d", http.StatusOK, w.Code)
        }

        // 验证响应体
        // ...
    })

    t.Run("验证查询参数", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/markets?limit=invalid", nil)
        w := httptest.NewRecorder()

        GetMarkets(w, req)

        if w.Code != http.StatusBadRequest {
            t.Errorf("期望状态码 %d，得到 %d", http.StatusBadRequest, w.Code)
        }
    })

    t.Run("优雅地处理数据库错误", func(t *testing.T) {
        // 模拟数据库失败
        req := httptest.NewRequest("GET", "/api/markets", nil)
        w := httptest.NewRecorder()

        GetMarkets(w, req)

        // 测试错误处理
    })
}
```


## 测试文件组织

```
pkg/
├── utils/
│   ├── math.go
│   └── math_test.go                 # 单元测试
├── handlers/
│   ├── markets.go
│   └── markets_test.go              # 集成测试
└── models/
    ├── market.go
    └── market_test.go
```

## 模拟外部服务

### 数据库模拟（使用接口）
```go
package db

// 定义接口
type Database interface {
    Query(query string, args ...interface{}) ([]Market, error)
    Insert(market Market) error
}

// 模拟实现
type MockDatabase struct {
    Markets []Market
    Error   error
}

func (m *MockDatabase) Query(query string, args ...interface{}) ([]Market, error) {
    if m.Error != nil {
        return nil, m.Error
    }
    return m.Markets, nil
}

func (m *MockDatabase) Insert(market Market) error {
    return m.Error
}
```

### Redis 模拟
```go
package redis

type RedisClient interface {
    SearchByVector(query []float64) ([]Market, error)
    CheckHealth() (bool, error)
}

type MockRedisClient struct {
    Results []Market
    Error   error
}

func (m *MockRedisClient) SearchByVector(query []float64) ([]Market, error) {
    if m.Error != nil {
        return nil, m.Error
    }
    return m.Results, nil
}

func (m *MockRedisClient) CheckHealth() (bool, error) {
    return true, m.Error
}
```

### 外部 API 模拟
```go
package api

type OpenAIClient interface {
    GenerateEmbedding(text string) ([]float64, error)
}

type MockOpenAIClient struct {
    Embedding []float64
    Error     error
}

func (m *MockOpenAIClient) GenerateEmbedding(text string) ([]float64, error) {
    if m.Error != nil {
        return nil, m.Error
    }
    return m.Embedding, nil
}
```

## 测试覆盖率验证

### 运行覆盖率报告
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 覆盖率阈值
```bash
# 在 CI/CD 中检查覆盖率
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "覆盖率 $COVERAGE% 低于 80%"
    exit 1
fi
```

## 避免常见测试错误

### ❌ 错误：测试实现细节
```go
// 不要测试内部状态
func TestCounter(t *testing.T) {
    c := NewCounter()
    c.count = 5  // 直接访问私有字段
    if c.count != 5 {
        t.Error("测试失败")
    }
}
```

### ✅ 正确：测试公共接口
```go
// 测试公共方法
func TestCounter(t *testing.T) {
    c := NewCounter()
    c.Increment()
    if c.Value() != 1 {
        t.Errorf("期望 1，得到 %d", c.Value())
    }
}
```

### ❌ 错误：硬编码测试数据
```go
// 难以维护
func TestProcess(t *testing.T) {
    result := Process("hardcoded", "data", 123)
    if result != "expected" {
        t.Error("测试失败")
    }
}
```

### ✅ 正确：使用表驱动测试
```go
// 易于维护和扩展
func TestProcess(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"测试用例1", "input1", "output1"},
        {"测试用例2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Process(tt.input)
            if result != tt.expected {
                t.Errorf("期望 %s，得到 %s", tt.expected, result)
            }
        })
    }
}
```

### ❌ 错误：没有测试隔离
```go
// 测试相互依赖
func TestCreateUser(t *testing.T) {
    user := CreateUser("test")
    // 依赖全局状态
}

func TestUpdateUser(t *testing.T) {
    // 依赖于前一个测试
    user := GetUser("test")
    UpdateUser(user)
}
```

### ✅ 正确：独立的测试
```go
// 每个测试设置自己的数据
func TestCreateUser(t *testing.T) {
    user := CreateTestUser()
    err := CreateUser(user)
    if err != nil {
        t.Errorf("创建用户失败: %v", err)
    }
}

func TestUpdateUser(t *testing.T) {
    user := CreateTestUser()
    CreateUser(user)  // 独立设置
    err := UpdateUser(user)
    if err != nil {
        t.Errorf("更新用户失败: %v", err)
    }
}
```

## 持续测试

### 开发期间的监视模式
```bash
# 使用 gotestsum 或类似工具
gotestsum --watch ./...
# 或使用 entr
find . -name "*.go" | entr -r go test ./...
```

### 提交前钩子
```bash
# .git/hooks/pre-commit
#!/bin/sh
go test ./... && go vet ./... && go fmt ./...
```

### CI/CD 集成
```yaml
# GitHub Actions
- name: 运行测试
  run: go test -v -race -coverprofile=coverage.out ./...
- name: 上传覆盖率
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.out
```

## 最佳实践

1. **先编写测试** - 始终使用 TDD
2. **每个测试一个断言** - 专注于单一行为
3. **描述性测试名称** - 解释测试内容
4. **安排-执行-断言** - 清晰的测试结构
5. **模拟外部依赖** - 隔离单元测试
6. **测试边界情况** - 空值、未定义、空、大值
7. **测试错误路径** - 不仅仅是快乐路径
8. **保持测试快速** - 单元测试每个 < 50ms
9. **测试后清理** - 无副作用
10. **审查覆盖率报告** - 识别差距

## 成功指标

- 达到 80% 以上的代码覆盖率
- 所有测试通过（绿色）
- 没有跳过或禁用的测试
- 快速的测试执行（单元测试 < 30 秒）
- 测试在生产前捕获错误

---

**记住**：测试不是可选的。它们是安全网，能够实现自信的重构、快速开发和生产可靠性。
