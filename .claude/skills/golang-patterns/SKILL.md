---
name: golang-patterns
description: 后端架构模式、API设计、数据库优化和Go、Gin、GORM的服务器端最佳实践。
---

# 后端开发模式

使用Go构建可扩展服务器端应用程序的后端架构模式和最佳实践。

## 何时激活

- 设计REST或GraphQL API端点
- 实现仓库、服务或控制器层
- 优化数据库查询（N+1、索引、连接池）
- 添加缓存（Redis、内存、HTTP缓存头）
- 设置后台作业或异步处理
- 为API构建错误处理和验证
- 构建中间件（认证、日志、速率限制）

## API设计模式

### RESTful API结构

```go
// ✅ 基于资源的URL
GET    /api/markets                 // 列出资源
GET    /api/markets/:id             // 获取单个资源
POST   /api/markets                 // 创建资源
PUT    /api/markets/:id             // 替换资源
PATCH  /api/markets/:id             // 更新资源
DELETE /api/markets/:id             // 删除资源

// ✅ 用于过滤、排序、分页的查询参数
GET /api/markets?status=active&sort=volume&limit=20&offset=0
```

### 仓库模式

```go
// 抽象数据访问逻辑
type MarketRepository interface {
    FindAll(ctx context.Context, filters *MarketFilters) ([]*Market, error)
    FindByID(ctx context.Context, id string) (*Market, error)
    Create(ctx context.Context, data *CreateMarketDTO) (*Market, error)
    Update(ctx context.Context, id string, data *UpdateMarketDTO) (*Market, error)
    Delete(ctx context.Context, id string) error
}

type GormMarketRepository struct {
    db *gorm.DB
}

func NewGormMarketRepository(db *gorm.DB) *GormMarketRepository {
    return &GormMarketRepository{db: db}
}

func (r *GormMarketRepository) FindAll(ctx context.Context, filters *MarketFilters) ([]*Market, error) {
    var markets []*Market
    query := r.db.WithContext(ctx).Model(&Market{})

    if filters != nil {
        if filters.Status != "" {
            query = query.Where("status = ?", filters.Status)
        }
        if filters.Limit > 0 {
            query = query.Limit(filters.Limit)
        }
    }

    if err := query.Find(&markets).Error; err != nil {
        return nil, fmt.Errorf("failed to find markets: %w", err)
    }

    return markets, nil
}

func (r *GormMarketRepository) FindByID(ctx context.Context, id string) (*Market, error) {
    var market Market
    if err := r.db.WithContext(ctx).First(&market, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to find market: %w", err)
    }
    return &market, nil
}

// 其他方法...
```

### 服务层模式

```go
// 业务逻辑与数据访问分离
type MarketService struct {
    marketRepo MarketRepository
    vectorStore VectorStore
}

func NewMarketService(marketRepo MarketRepository, vectorStore VectorStore) *MarketService {
    return &MarketService{
        marketRepo: marketRepo,
        vectorStore: vectorStore,
    }
}

func (s *MarketService) SearchMarkets(ctx context.Context, query string, limit int) ([]*Market, error) {
    // 业务逻辑
    embedding, err := s.generateEmbedding(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to generate embedding: %w", err)
    }

    results, err := s.vectorSearch(ctx, embedding, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to vector search: %w", err)
    }

    // 获取完整数据
    ids := make([]string, len(results))
    for i, r := range results {
        ids[i] = r.ID
    }

    markets, err := s.marketRepo.FindByIDs(ctx, ids)
    if err != nil {
        return nil, fmt.Errorf("failed to find markets: %w", err)
    }

    // 按相似度排序
    sort.Slice(markets, func(i, j int) bool {
        scoreI := s.getScore(results, markets[i].ID)
        scoreJ := s.getScore(results, markets[j].ID)
        return scoreI > scoreJ
    })

    return markets, nil
}

func (s *MarketService) vectorSearch(ctx context.Context, embedding []float64, limit int) ([]*SearchResult, error) {
    // 向量搜索实现
    return s.vectorStore.Search(ctx, embedding, limit)
}

func (s *MarketService) getScore(results []*SearchResult, id string) float64 {
    for _, r := range results {
        if r.ID == id {
            return r.Score
        }
    }
    return 0
}
```

### 中间件模式

```go
// 请求/响应处理管道
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            c.Abort()
            return
        }

        user, err := verifyToken(token, jwtSecret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user", user)
        c.Next()
    }
}

// 使用
func SetupRoutes(r *gin.Engine, jwtSecret string) {
    api := r.Group("/api")
    api.Use(AuthMiddleware(jwtSecret))
    {
        api.GET("/markets", GetMarkets)
        api.POST("/markets", CreateMarket)
    }
}

func GetMarkets(c *gin.Context) {
    // 处理程序可以通过c.Get("user")访问用户
    user, _ := c.Get("user")
    // 处理请求...
}
```

## 数据库模式

### 查询优化

```go
// ✅ GOOD: 只选择需要的列
var markets []Market
err := db.
    Select("id", "name", "status", "volume").
    Where("status = ?", "active").
    Order("volume DESC").
    Limit(10).
    Find(&markets).Error

// ❌ BAD: 选择所有内容
var markets []Market
err := db.Find(&markets).Error
```

### N+1查询预防

```go
// ❌ BAD: N+1查询问题
var markets []Market
db.Find(&markets)

for _, market := range markets {
    var user User
    db.First(&user, market.CreatorID)  // N个查询
    market.Creator = &user
}
}

// ✅ GOOD: 使用Preload批量获取
var markets []Market
db.
    Preload("Creator").  // 1个JOIN查询
    Find(&markets)

// ✅ GOOD: 手动批量获取
var markets []Market
db.Find(&markets)

creatorIDs := make([]uint, len(markets))
for i, m := range markets {
    creatorIDs[i] = m.CreatorID
}

var creators []User
db.Where("id IN ?", creatorIDs).Find(&creators)

creatorMap := make(map[uint]*User)
for i := range creators {
    creatorMap[creators[i].ID] = &creators[i]
}

for i := range markets {
    markets[i].Creator = creatorMap[markets[i].CreatorID]
}
```

### 事务模式

```go
func CreateMarketWithPosition(
    ctx context.Context,
    db *gorm.DB,
    marketData *CreateMarketDTO,
    positionData *CreatePositionDTO,
) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 创建市场
        market := &Market{
            Name:   marketData.Name,
            Status: marketData.Status,
        }
        if err := tx.Create(market).Error; err != nil {
            return fmt.Errorf("failed to create market: %w", err)
        }

        // 创建位置
        position := &Position{
            MarketID: market.ID,
            Amount:   positionData.Amount,
        }
        if err := tx.Create(position).Error; err != nil {
            return fmt.Errorf("failed to create position: %w", err)
        }

        return nil
    })
}

// 使用
err := CreateMarketWithPosition(ctx, db, marketData, positionData)
if err != nil {
    // 事务自动回滚
    log.Printf("Transaction failed: %v", err)
}
```

## 缓存策略

### Redis缓存层

```go
type CachedMarketRepository struct {
    baseRepo MarketRepository
    redis    *redis.Client
}

func NewCachedMarketRepository(baseRepo MarketRepository, redis *redis.Client) *CachedMarketRepository {
    return &CachedMarketRepository{
        baseRepo: baseRepo,
        redis:    redis,
    }
}

func (r *CachedMarketRepository) FindByID(ctx context.Context, id string) (*Market, error) {
    // 首先检查缓存
    cacheKey := fmt.Sprintf("market:%s", id)
    cached, err := r.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var market Market
        if err := json.Unmarshal([]byte(cached), &market); err == nil {
            return &market, nil
        }
    }

    // 缓存未命中 - 从数据库获取
    market, err := r.baseRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if market != nil {
        // 缓存5分钟
        data, _ := json.Marshal(market)
        r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
    }

    return market, nil
}

func (r *CachedMarketRepository) InvalidateCache(ctx context.Context, id string) error {
    cacheKey := fmt.Sprintf("market:%s", id)
    return r.redis.Del(ctx, cacheKey).Err()
}
```

### 缓存旁路模式

```go
func GetMarketWithCache(ctx context.Context, db *gorm.DB, redis *redis.Client, id string) (*Market, error) {
    cacheKey := fmt.Sprintf("market:%s", id)

    // 尝试缓存
    cached, err := redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var market Market
        if err := json.Unmarshal([]byte(cached), &market); err == nil {
            return &market, nil
        }
    }

    // 缓存未命中 - 从数据库获取
    var market Market
    if err := db.First(&market, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("market not found")
        }
        return nil, fmt.Errorf("failed to fetch market: %w", err)
    }

    // 更新缓存
    data, _ := json.Marshal(market)
    redis.Set(ctx, cacheKey, data, 5*time.Minute)

    return &market, nil
}
```

## 错误处理模式

### 集中错误处理程序

```go
type AppError struct {
    StatusCode int
    Message    string
    IsOperational bool
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(statusCode int, message string, isOperational bool) *AppError {
    return &AppError{
        StatusCode: statusCode,
        Message:    message,
        IsOperational: isOperational,
    }
}

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            switch e := err.(type) {
            case *AppError:
                c.JSON(e.StatusCode, gin.H{
                    "success": false,
                    "error":   e.Message,
                })
            case *validator.ValidationErrors:
                c.JSON(http.StatusBadRequest, gin.H{
                    "success": false,
                    "error":   "Validation failed",
                    "details": formatValidationErrors(*e),
                })
            default:
                // 记录意外错误
                log.Printf("Unexpected error: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "error":   "Internal server error",
                })
            }
        }
    }
}

func formatValidationErrors(errs validator.ValidationErrors) []map[string]string {
    details := make([]map[string]string, len(errs))
    for i, e := range errs {
        details[i] = map[string]string{
            "field":   e.Field(),
            "message": fmt.Sprintf("failed on '%s' tag", e.Tag()),
        }
    }
    return details
}

// 使用
func GetMarket(c *gin.Context) {
    id := c.Param("id")

    market, err := marketService.GetByID(c.Request.Context(), id)
    if err != nil {
        c.Error(err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": market})
}
```

### 指数退避重试

```go
func RetryWithBackoff[T any](
    ctx context.Context,
    fn func() (T, error),
    maxRetries int,
) (T, error) {
    var result T
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        result, lastErr = fn()
        if lastErr == nil {
            return result, nil
        }

        if i < maxRetries-1 {
            // 指数退避: 1s, 2s, 4s
            delay := time.Duration(math.Pow(2, float64(i))) * time.Second
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return result, ctx.Err()
            }
        }
    }

    return result, lastErr
}

// 使用
data, err := RetryWithBackoff(ctx, func() ([]byte, error) {
    return fetchFromAPI()
}, 3)
```

## 认证与授权

### JWT令牌验证

```go
type JWTPayload struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
}

func GenerateToken(payload *JWTPayload, secret string, expiration time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "user_id": payload.UserID,
        "email":   payload.Email,
        "role":    payload.Role,
        "exp":     time.Now().Add(expiration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString, secret string) (*JWTPayload, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, NewAppError(http.StatusUnauthorized, "Invalid token", true)
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return &JWTPayload{
            UserID: claims["user_id"].(string),
            Email:  claims["email"].(string),
            Role:   claims["role"].(string),
        }, nil
    }

    return nil, NewAppError(http.StatusUnauthorized, "Invalid token", true)
}

func RequireAuth(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.Error(NewAppError(http.StatusUnauthorized, "Missing authorization token", true))
        c.Abort()
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        c.Error(NewAppError(http.StatusUnauthorized, "Invalid authorization header", true))
        c.Abort()
        return
    }

    secret := c.GetString("jwt_secret")
    user, err := VerifyToken(token, secret)
    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.Set("user", user)
    c.Next()
}

// 在处理程序中使用
func GetUserData(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    payload := user.(*JWTPayload)
    data, err := getDataForUser(c.Request.Context(), payload.UserID)
    if err != nil {
        c.Error(err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
```

## 速率限制

### 简单的内存速率限制器

```go
type RateLimiter struct {
    requests map[string][]int64
    mu       sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]int64),
    }
}

func (rl *RateLimiter) CheckLimit(identifier string, maxRequests int, windowMs int64) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now().UnixMilli()
    requests := rl.requests[identifier]

    // 移除窗口外的旧请求
    var recentRequests []int64
    for _, t := range requests {
        if now-t < windowMs {
            recentRequests = append(recentRequests, t)
        }
    }

    if len(recentRequests) >= maxRequests {
        return false // 超出速率限制
    }

    // 添加当前请求
    recentRequests = append(recentRequests, now)
    rl.requests[identifier] = recentRequests

    return true
}

var limiter = NewRateLimiter()

func RateLimitMiddleware(maxRequests int, windowMs int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        allowed := limiter.CheckLimit(ip, maxRequests, windowMs)
        if !allowed {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用
func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api")
    api.Use(RateLimitMiddleware(100, 60000)) // 100 req/min
    {
        api.GET("/markets", GetMarkets)
    }
}
```

## 后台作业和队列

### 简单队列模式

```go
type Job interface {
    Execute(ctx context.Context) error
}

type JobQueue struct {
    queue     chan Job
    workers   int
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewJobQueue(workers int, queueSize int) *JobQueue {
    ctx, cancel := context.WithCancel(context.Background())
    return &JobQueue{
        queue:   make(chan Job, queueSize),
        workers: workers,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (jq *JobQueue) Start() {
    for i := 0; i < jq.workers; i++ {
        jq.wg.Add(1)
        go jq.worker()
    }
}

func (jq *JobQueue) worker() {
    defer jq.wg.Done()

    for {
        select {
        case job := <-jq.queue:
            if err := job.Execute(jq.ctx); err != nil {
                log.Printf("Job failed: %v", err)
            }
        case <-jq.ctx.Done():
            return
        }
    }
}

func (jq *JobQueue) Add(job Job) error {
    select {
    case jq.queue <- job:
        return nil
    case <-jq.ctx.Done():
        return fmt.Errorf("queue is shutting down")
    }
}

func (jq *JobQueue) Stop() {
    jq.cancel()
    jq.wg.Wait()
}

// 用于索引市场的使用示例
type IndexJob struct {
    MarketID string
    Service  *MarketService
}

func (j *IndexJob) Execute(ctx context.Context) error {
    return j.Service.IndexMarket(ctx, j.MarketID)
}

var indexQueue = NewJobQueue(5, 1000)

func init() {
    indexQueue.Start()
}

func CreateMarket(c *gin.Context) {
    var req CreateMarketRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    market, err := marketService.Create(c.Request.Context(), &req)
    if err != nil {
        c.Error(err)
        return
    }

    // 添加到队列而不是阻塞
    job := &IndexJob{MarketID: market.ID, Service: marketService}
    if err := indexQueue.Add(job); err != nil {
        log.Printf("Failed to queue index job: %v", err)
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Market created and indexing queued",
        "data":    market,
    })
}
```

## 日志和监控

### 结构化日志

```go
type LogContext struct {
    UserID    string `json:"user_id,omitempty"`
    RequestID string `json:"request_id,omitempty"`
    Method    string `json:"method,omitempty"`
    Path      string `json:"path,omitempty"`
    // 根据需要添加更多字段
}

type Logger struct {
    logger *logrus.Logger
}

func NewLogger() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetLevel(logrus.InfoLevel)
    return &Logger{logger: logger}
}

func (l *Logger) log(level logrus.Level, message string, context *LogContext) {
    entry := l.logger.WithFields(logrus.Fields{
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "level":     level.String(),
        "message":   message,
    })

    if context != nil {
        entry = entry.WithFields(logrus.Fields{
            "user_id":    context.UserID,
            "request_id": context.RequestID,
            "method":     context.Method,
            "path":       context.Path,
        })
    }

    entry.Log(level)
}

func (l *Logger) Info(message string, context *LogContext) {
    l.log(logrus.InfoLevel, message, context)
}

func (l *Logger) Warn(message string, context *LogContext) {
    l.log(logrus.WarnLevel, message, context)
}

func (l *Logger) Error(message string, err error, context *LogContext) {
    entry := l.logger.WithFields(logrus.Fields{
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "level":     "error",
        "message":   message,
        "error":     err.Error(),
    })

    if context != nil {
        entry = entry.WithFields(logrus.Fields{
            "user_id":    context.UserID,
            "request_id": context.RequestID,
            "method":     context.Method,
            "path":       context.Path,
        })
    }

    entry.Error(err)
}

var logger = NewLogger()

// 使用
func GetMarkets(c *gin.Context) {
    requestID := uuid.New().String()

    logger.Info("Fetching markets", &LogContext{
        RequestID: requestID,
        Method:    c.Request.Method,
        Path:      c.Request.URL.Path,
    })

    markets, err := marketService.FindAll(c.Request.Context())
    if err != nil {
        logger.Error("Failed to fetch markets", err, &LogContext{
            RequestID: requestID,
        })
        c.Error(err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": markets})
}
```

### 请求日志中间件

```go
func RequestLoggerMiddleware(logger *Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)

        // 记录请求
        logger.Info("Request started", &LogContext{
            RequestID: requestID,
            Method:    c.Request.Method,
            Path:      c.Request.URL.Path,
        })

        c.Next()

        // 记录响应
        duration := time.Since(start)
        logger.Info("Request completed", &LogContext{
            RequestID: requestID,
            Method:    c.Request.Method,
            Path:      c.Request.URL.Path,
        })
    }
}
```

## 项目结构

### 推荐的Go项目布局

```
project/
├── cmd/
│   └── server/
│       └── main.go              // 应用程序入口点
├── internal/
│   ├── config/
│   │   └── config.go            // 配置管理
│   ├── domain/
│   │   ├── market.go            // 领域模型
│   │   └── user.go
│   ├── repository/
│   │   ├── market.go            // 数据访问层
│   │   └── user.go
│   ├── service/
│   │   ├── market.go            // 业务逻辑层
│   │   └── user.go
│   ├── handler/
│   │   ├── market.go            // HTTP处理程序
│   │   └── user.go
│   └── middleware/
│       ├── auth.go
│       ├── logging.go
│       └── ratelimit.go
├── pkg/
│   ├── cache/
│   │   └── redis.go             // Redis工具
│   └── errors/
│       └── errors.go            // 错误定义
├── migrations/                  // 数据库迁移
├── go.mod
├── go.sum
└── README.md
```

### 依赖注入示例

```go
type Container struct {
    Config     *config.Config
    DB         *gorm.DB
    Redis      *redis.Client
    Logger     *Logger
    
    // 仓库
    MarketRepo repository.MarketRepository
    UserRepo   repository.UserRepository
    
    // 服务
    MarketService *service.MarketService
    UserService   *service.UserService
    
    // 处理程序
    MarketHandler *handler.MarketHandler
    UserHandler   *handler.UserHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
    c := &Container{Config: cfg}
    
    // 初始化数据库
    db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    c.DB = db
    
    // 初始化Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.Addr,
    })
    c.Redis = rdb
    
    // 初始化日志记录器
    c.Logger = NewLogger()
    
    // 初始化仓库
    c.MarketRepo = repository.NewGormMarketRepository(db)
    c.UserRepo = repository.NewGormUserRepository(db)
    
    // 初始化服务
    c.MarketService = service.NewMarketService(c.MarketRepo, c.Redis)
    c.UserService = service.NewUserService(c.UserRepo)
    
    // 初始化处理程序
    c.MarketHandler = handler.NewMarketHandler(c.MarketService)
    c.UserHandler = handler.NewUserHandler(c.UserService)
    
    return c, nil
}
```

**记住**: 后端模式能够构建可扩展、可维护的服务器端应用程序。选择适合你复杂度级别的模式。Go的简洁性和强类型使其成为构建强大后端服务的理想选择。