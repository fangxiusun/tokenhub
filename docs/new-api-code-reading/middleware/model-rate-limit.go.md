# model-rate-limit.go 代码阅读文档

## 1. 全局总结
model-rate-limit.go 实现了每模型请求限流中间件，允许管理员为特定模型设置请求频率限制。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `dto` - 数据传输对象
- `setting` - 设置管理
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/relay-router.go` - 中转路由

## 3. 类型定义
### 3.1 函数
- `ModelRequestRateLimit() gin.HandlerFunc` - 模型请求限流中间件

## 4. 函数详解
### 4.1 ModelRequestRateLimit
- **签名**: `func ModelRequestRateLimit() gin.HandlerFunc`
- **职责**: 按模型名称限制请求频率
- **逻辑流程**:
  1. 检查是否启用模型限流
  2. 解析请求体获取模型名称
  3. 查询模型限流配置
  4. 检查 IP 是否超限
  5. 如果超限返回 429

## 5. 关键逻辑分析
- **按模型配置**: 每个模型可独立配置限制
- **分组支持**: 可按用户组配置不同限制
- **动态配置**: 通过管理面板实时更新限制

## 6. 关联文件
- `setting/rate_limit.go` - 模型限流配置
- `common/rate-limit.go` - 内存限流器
