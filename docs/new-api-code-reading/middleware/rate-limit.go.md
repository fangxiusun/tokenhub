# rate-limit.go 代码阅读文档

## 1. 全局总结
rate-limit.go 实现了全局限流中间件，支持 Redis 和内存两种后端。使用滑动窗口算法，防止 API 被滥用。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具（Redis、内存限流器）
- `constant` - 常量定义
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/api-router.go` - 管理 API 路由
- `router/relay-router.go` - 中转路由

## 3. 类型定义
### 3.1 函数
- `GlobalAPIRateLimit() gin.HandlerFunc` - 全局 API 限流中间件
- `CriticalRateLimit() gin.HandlerFunc` - 关键操作限流（登录、注册等）

## 4. 函数详解
### 4.1 GlobalAPIRateLimit
- **签名**: `func GlobalAPIRateLimit() gin.HandlerFunc`
- **职责**: 按 IP 限制 API 请求频率
- **逻辑流程**:
  1. 获取客户端 IP
  2. 检查 Redis 或内存限流器
  3. 如果超限返回 429 状态码
  4. 记录请求到限流器

### 4.2 CriticalRateLimit
- **签名**: `func CriticalRateLimit() gin.HandlerFunc`
- **职责**: 限制关键操作频率（20 次/20 分钟）
- **逻辑流程**: 类似 GlobalAPIRateLimit，但使用更严格的限制

## 5. 关键逻辑分析
- **滑动窗口算法**: 使用时间戳记录请求，清理过期记录
- **双后端支持**: Redis 优先，不可用时回退到内存
- **IP 提取**: 支持 X-Forwarded-For 和 X-Real-IP 头

## 6. 关联文件
- `common/redis.go` - Redis 限流器
- `common/rate-limit.go` - 内存限流器
- `common/limiter/limiter.go` - Redis Lua 脚本限流
