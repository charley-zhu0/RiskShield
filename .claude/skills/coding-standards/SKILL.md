---
name: coding-standards
description: Universal coding standards, best practices, and patterns for Go development.
---

# Go 编码规范与最佳实践

适用于 Go 项目的通用编码规范。

## When to Activate

- 启动新的 Go 项目或模块
- 审查代码质量和可维护性
- 重构现有代码以遵循规范
- 强制命名、格式或结构一致性
- 设置 linting、格式化或静态检查规则
- 新贡献者的编码规范培训

## 代码质量原则

### 1. 可读性优先
- 代码被阅读的次数远多于编写
- 清晰的变量和函数命名
- 自解释代码优于注释
- 一致的格式化 (`gofmt`/`goimports`)

### 2. KISS (保持简单)
- 使用最简单的可行方案
- 避免过度工程化
- 不做过早优化
- 易于理解 > 聪明的代码

### 3. DRY (不要重复自己)
- 将通用逻辑抽取为函数
- 创建可复用的包
- 跨模块共享工具函数
- 避免复制粘贴编程

### 4. YAGNI (你不会需要它)
- 不要在需要之前构建功能
- 避免投机性的泛化
- 仅在需要时增加复杂性
- 从简单开始，按需重构

## Go 命名规范

### 包命名

```go
// ✅ 推荐: 简短、小写、无下划线
package user
package httputil
package testdata

// ❌ 避免: 混合大小写、下划线、通用名称
package userService   // 不要用驼峰
package http_util     // 不要用下划线
package util          // 太通用
package common        // 太通用
package base          // 太通用
```

### 变量命名

```go
// ✅ 推荐: 驼峰式，根据作用域选择长度
userID := getUserID()              // 局部变量简短
isAuthenticated := checkAuth()     // 布尔值用 is/has/can 前缀
maxRetryCount := 3                 // 常量描述清晰

// 循环变量
for i, user := range users { }     // 短作用域用单字母
for idx, item := range items { }   // 需要更多上下文时稍长

// ❌ 避免: 不清晰或过长的名称
u := getUser()                     // 太短，不清晰
theCurrentUserObjectInstance := x  // 太长
```

### 函数命名

```go
// ✅ 推荐: 动词+名词模式，导出函数首字母大写
func GetUser(id string) (*User, error) { }
func (s *Service) CreateOrder(ctx context.Context, req *CreateOrderRequest) error { }
func parseConfig(data []byte) (*Config, error) { }  // 未导出
func isValidEmail(email string) bool { }            // 返回 bool 用 is/has/can

// 构造函数
func NewUserService(db *sql.DB) *UserService { }    // New + 类型名
func NewUserServiceWithOptions(opts ...Option) *UserService { }

// ❌ 避免
func user(id string) { }           // 缺少动词
func DoGetUserData(id string) { }  // 冗余前缀
```

### 接口命名

```go
// ✅ 推荐: 单方法接口用 -er 后缀
type Reader interface {
    Read(p []byte) (n int, err error)
}

type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// ❌ 避免: I 前缀（Java 风格）
type IUserRepository interface { }  // 不要用 I 前缀
```

### 常量和枚举

```go
// ✅ 推荐: 使用 iota 定义枚举
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusClosed
)

// 字符串常量组
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
    BufferSize     = 4096
)

// ❌ 避免: 无类型常量用于枚举场景
const (
    PENDING = 0  // 全大写是其他语言风格
    ACTIVE  = 1
)
```

## 错误处理

### 基本错误处理

```go
// ✅ 推荐: 立即处理错误，提供上下文
user, err := s.repo.GetByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("get user %s: %w", userID, err)
}

// 使用 errors.Is 和 errors.As 检查错误
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrUserNotFound
}

var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == "23505" {
    return nil, ErrDuplicateEmail
}

// ❌ 避免: 忽略错误或丢失上下文
user, _ := s.repo.GetByID(ctx, userID)  // 忽略错误
return nil, err                          // 丢失上下文
```

### 自定义错误

```go
// ✅ 推荐: 定义领域特定的哨兵错误
var (
    ErrUserNotFound   = errors.New("user not found")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized   = errors.New("unauthorized")
)

// 包含更多信息的错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// 检查自定义错误
var validErr *ValidationError
if errors.As(err, &validErr) {
    log.Printf("validation failed: %s", validErr.Field)
}
```

### defer 中的错误处理

```go
// ✅ 推荐: 在 defer 中处理可能返回错误的关闭操作
func processFile(path string) (err error) {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open file: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("close file: %w", cerr)
        }
    }()
    
    // 处理文件...
    return nil
}
```

## 并发编程

### Goroutine 管理

```go
// ✅ 推荐: 使用 errgroup 管理 goroutine 生命周期
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]Result, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]Result, len(urls))
    
    for i, url := range urls {
        i, url := i, url  // Go 1.22 之前需要捕获循环变量
        g.Go(func() error {
            result, err := fetch(ctx, url)
            if err != nil {
                return err
            }
            results[i] = result
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}

// ❌ 避免: 裸 goroutine 无错误处理
for _, url := range urls {
    go fetch(url)  // 错误丢失，无法等待完成
}
```

### Channel 使用

```go
// ✅ 推荐: 清晰的 channel 所有权
func producer(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)  // 生产者负责关闭
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                return
            case ch <- i:
            }
        }
    }()
    return ch
}

// 使用 select 处理超时
select {
case result := <-resultCh:
    return result, nil
case <-ctx.Done():
    return nil, ctx.Err()
case <-time.After(5 * time.Second):
    return nil, ErrTimeout
}
```

### 互斥锁

```go
// ✅ 推荐: 互斥锁紧贴被保护的数据
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    return item, ok
}

func (c *Cache) Set(key string, item Item) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = item
}

// ❌ 避免: 锁作用域过大
func (c *Cache) Process(key string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    item := c.items[key]
    result := expensiveOperation(item)  // 不需要锁
    c.items[key] = result
    return nil
}
```

## Context 使用

### 传递 Context

```go
// ✅ 推荐: context 作为第一个参数
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // 检查 context 是否已取消
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    return s.repo.GetByID(ctx, id)
}

// 在 HTTP handler 中
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    user, err := h.service.GetUser(ctx, chi.URLParam(r, "id"))
    // ...
}

// ❌ 避免
func GetUser(id string, ctx context.Context) { }  // context 不是第一个参数
func GetUser(id string) { }                        // 缺少 context
```

### Context 值

```go
// ✅ 推荐: 使用自定义类型作为 key 避免冲突
type contextKey string

const (
    userIDKey   contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}

// ❌ 避免: 使用字符串或内置类型作为 key
ctx = context.WithValue(ctx, "userID", userID)  // 可能冲突
```

## 结构体设计

### 结构体定义

```go
// ✅ 推荐: 字段按逻辑分组，考虑内存对齐
type User struct {
    // 标识符
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 基本信息
    Name  string `json:"name"`
    Email string `json:"email"`
    
    // 状态
    IsActive bool   `json:"is_active"`
    Role     string `json:"role"`
}

// 使用 functional options 模式
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *slog.Logger) ServerOption {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(addr string, opts ...ServerOption) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second,  // 默认值
        logger:  slog.Default(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### 方法接收者

```go
// ✅ 推荐: 一致的接收者命名，通常 1-2 个字母
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

func (s *UserService) Create(ctx context.Context, u *User) error {
    // ...
}

// 值接收者 vs 指针接收者
// 使用指针: 需要修改接收者，大型结构体，或为保持一致性
func (u *User) SetName(name string) {
    u.Name = name
}

// 使用值: 小型不可变结构体，基本类型包装
type Point struct{ X, Y float64 }
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// ❌ 避免: 不一致的接收者命名
func (user *User) GetName() string { }     // 太长
func (this *User) SetEmail(e string) { }   // 不要用 this/self
```

## 项目结构

### 标准项目布局

```
myproject/
├── cmd/                    # 主应用程序入口
│   ├── api/               # API 服务器
│   │   └── main.go
│   └── worker/            # 后台 worker
│       └── main.go
├── internal/              # 私有应用代码（不可被外部导入）
│   ├── config/           # 配置加载
│   ├── domain/           # 领域模型
│   ├── handler/          # HTTP handlers
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                   # 可被外部项目导入的库代码
│   └── httputil/
├── api/                   # API 定义文件
│   └── openapi.yaml
├── migrations/            # 数据库迁移文件
├── scripts/               # 构建和工具脚本
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 文件命名

```
user.go                    # 主要类型定义
user_repository.go         # 仓库实现
user_service.go            # 服务层
user_handler.go            # HTTP handler
user_test.go               # 测试文件
user_mock.go               # Mock 实现
```

## 测试规范

### 表驱动测试

```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        price    float64
        quantity int
        want     float64
        wantErr  bool
    }{
        {
            name:     "no discount for small order",
            price:    100,
            quantity: 5,
            want:     500,
        },
        {
            name:     "10% discount for orders over 10",
            price:    100,
            quantity: 15,
            want:     1350,
        },
        {
            name:     "error for negative price",
            price:    -100,
            quantity: 1,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CalculateDiscount(tt.price, tt.quantity)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got = %v, want = %v", got, tt.want)
            }
        })
    }
}
```

### 使用 testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserService_Create(t *testing.T) {
    // Arrange
    repo := NewMockRepository()
    service := NewUserService(repo)
    user := &User{Name: "Alice", Email: "alice@example.com"}
    
    // Act
    err := service.Create(context.Background(), user)
    
    // Assert
    require.NoError(t, err)  // 失败则立即停止
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, "Alice", user.Name)
}
```

### 子测试和并行测试

```go
func TestUser(t *testing.T) {
    t.Run("Create", func(t *testing.T) {
        t.Parallel()  // 可并行执行
        // ...
    })
    
    t.Run("Update", func(t *testing.T) {
        t.Parallel()
        // ...
    })
}
```

## 日志规范

### 使用 slog (Go 1.21+)

```go
import "log/slog"

// ✅ 推荐: 结构化日志
logger := slog.Default()

logger.Info("user created",
    slog.String("user_id", user.ID),
    slog.String("email", user.Email),
)

logger.Error("failed to create user",
    slog.String("email", email),
    slog.Any("error", err),
)

// 带 context 的日志
logger.InfoContext(ctx, "processing request",
    slog.String("request_id", requestID),
)

// 创建子 logger
userLogger := logger.With(slog.String("user_id", userID))
userLogger.Info("user action")

// ❌ 避免
log.Printf("user created: %s", user.ID)  // 非结构化
fmt.Println("debug:", data)              // 不要用 fmt 做日志
```

## 性能最佳实践

### 避免不必要的内存分配

```go
// ✅ 推荐: 预分配切片容量
users := make([]User, 0, len(ids))
for _, id := range ids {
    user, _ := getUser(id)
    users = append(users, user)
}

// 使用 strings.Builder 拼接字符串
var b strings.Builder
b.Grow(100)  // 预估容量
for _, s := range parts {
    b.WriteString(s)
}
result := b.String()

// ❌ 避免
var result string
for _, s := range parts {
    result += s  // 每次都分配新内存
}
```

### 使用 sync.Pool

```go
var bufferPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func process(data []byte) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    buf.Write(data)
    // 使用 buffer...
}
```

### 避免闭包捕获问题 (Go < 1.22)

```go
// ✅ 推荐: 显式传递循环变量 (Go < 1.22)
for _, item := range items {
    item := item  // 创建副本
    go func() {
        process(item)
    }()
}

// 或者作为参数传递
for _, item := range items {
    go func(i Item) {
        process(i)
    }(item)
}

// Go 1.22+ 已自动修复此问题，循环变量每次迭代都是新变量
```

## 代码异味检测

关注以下反模式:

### 1. 过长函数
```go
// ❌ 避免: 函数超过 50 行
func processOrder(order *Order) error {
    // 100 行代码...
}

// ✅ 推荐: 拆分为小函数
func processOrder(order *Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    if err := calculateTotal(order); err != nil {
        return err
    }
    return saveOrder(order)
}
```

### 2. 过深嵌套
```go
// ❌ 避免: 5+ 层嵌套
if user != nil {
    if user.IsActive {
        if order != nil {
            if order.IsValid {
                // 做某事
            }
        }
    }
}

// ✅ 推荐: 提前返回
if user == nil {
    return ErrUserNotFound
}
if !user.IsActive {
    return ErrUserInactive
}
if order == nil {
    return ErrOrderNotFound
}
if !order.IsValid {
    return ErrInvalidOrder
}
// 做某事
```

### 3. 过大的接口
```go
// ❌ 避免: 接口太大
type Repository interface {
    GetUser(id string) (*User, error)
    CreateUser(u *User) error
    UpdateUser(u *User) error
    DeleteUser(id string) error
    GetOrder(id string) (*Order, error)
    CreateOrder(o *Order) error
    // ... 20+ 方法
}

// ✅ 推荐: 小接口，单一职责
type UserRepository interface {
    Get(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, u *User) error
    Update(ctx context.Context, u *User) error
    Delete(ctx context.Context, id string) error
}

type OrderRepository interface {
    Get(ctx context.Context, id string) (*Order, error)
    Create(ctx context.Context, o *Order) error
}
```

## 工具配置

### golangci-lint 配置

```yaml
# .golangci.yml
run:
  timeout: 5m
  go: '1.21'

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - gocritic
    - revive
    - errorlint
    - exhaustive

linters-settings:
  govet:
    enable-all: true
  revive:
    rules:
      - name: blank-imports
      - name: context-as-argument
      - name: error-return
      - name: error-naming
      - name: exported
      - name: var-naming

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
        - gocritic
```

### Makefile 示例

```makefile
.PHONY: build test lint fmt

GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w"

build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/api ./cmd/api

test:
	$(GO) test $(GOFLAGS) -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	goimports -w .

tidy:
	$(GO) mod tidy
	$(GO) mod verify