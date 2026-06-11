# cors.go 代码阅读文档

## 1. 全局总结
cors.go 实现了 CORS（跨域资源共享）中间件，允许前端跨域访问 API。

## 2. 依赖关系
### 2.1 导入的包
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/main.go` - 路由初始化

## 3. 类型定义
### 3.1 函数
- `Cors() gin.HandlerFunc` - CORS 中间件

## 4. 函数详解
### 4.1 Cors
- **签名**: `func Cors() gin.HandlerFunc`
- **职责**: 设置 CORS 响应头
- **逻辑流程**:
  1. 设置 Access-Control-Allow-Origin: *
  2. 设置 Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
  3. 设置 Access-Control-Allow-Headers: Content-Type, Authorization
  4. 处理 OPTIONS 预检请求

## 5. 关键逻辑分析
- **全开放策略**: 允许所有来源访问（生产环境建议限制）
- **预检请求**: OPTIONS 请求直接返回 200

## 6. 关联文件
- `router/main.go` - 路由初始化中应用此中间件
