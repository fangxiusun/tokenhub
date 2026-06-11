# logger.go 代码阅读文档

## 1. 全局总结
logger.go 实现了请求日志中间件，记录每个 HTTP 请求的详细信息。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/main.go` - 路由初始化

## 3. 类型定义
### 3.1 函数
- `Logger() gin.HandlerFunc` - 请求日志中间件

## 4. 函数详解
### 4.1 Logger
- **签名**: `func Logger() gin.HandlerFunc`
- **职责**: 记录请求方法、路径、状态码、耗时
- **逻辑流程**:
  1. 记录请求开始时间
  2. 处理请求
  3. 计算耗时
  4. 记录日志（带请求 ID）

## 5. 关键逻辑分析
- **请求 ID 关联**: 日志包含请求 ID，便于追踪
- **耗时统计**: 计算请求处理时间
- **错误日志**: 非 2xx 状态码记录为错误

## 6. 关联文件
- `common/sys_log.go` - 系统日志
- `middleware/request-id.go` - 请求 ID 生成
