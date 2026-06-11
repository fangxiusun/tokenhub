# gzip.go 代码阅读文档

## 1. 全局总结
gzip.go 实现了 Gzip 压缩中间件，减少响应体大小，提高传输效率。

## 2. 依赖关系
### 2.1 导入的包
- `github.com/gin-contrib/gzip` - Gzip 压缩库
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/web-router.go` - 静态资源路由

## 3. 类型定义
### 3.1 函数
- `Gzip(level int) gin.HandlerFunc` - Gzip 压缩中间件

## 4. 函数详解
### 4.1 Gzip
- **签名**: `func Gzip(level int) gin.HandlerFunc`
- **职责**: 对响应体进行 Gzip 压缩
- **参数**: level - 压缩级别（-1 到 9）

## 5. 关键逻辑分析
- **压缩级别**: 使用 gzip.DefaultCompression 作为默认级别
- **Content-Type 检查**: 只压缩文本类型的响应

## 6. 关联文件
- `router/web-router.go` - 静态资源路由中应用此中间件
