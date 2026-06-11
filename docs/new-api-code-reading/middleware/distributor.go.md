# distributor.go 代码阅读文档

## 1. 全局总结
distributor.go 实现了渠道路由分发中间件，是系统的核心路由组件。它根据请求的模型名称和用户组，从可用渠道中选择一个合适的渠道，并将渠道信息注入到请求上下文中。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `dto` - 数据传输对象
- `model` - 数据模型
- `service` - 业务逻辑（渠道选择）
- `types` - 类型定义
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/relay-router.go` - 中转路由
- `controller/relay.go` - 中转控制器

## 3. 类型定义
### 3.1 函数
- `Distribute() gin.HandlerFunc` - 渠道路由分发中间件
- `parseBody(c *gin.Context) map[string]interface{}` - 解析请求体
- `getModelName(c *gin.Context, body map[string]interface{}) string` - 获取模型名称

## 4. 函数详解
### 4.1 Distribute
- **签名**: `func Distribute() gin.HandlerFunc`
- **职责**: 根据模型和用户组选择渠道
- **逻辑流程**:
  1. 检查是否启用渠道路由
  2. 解析请求体获取模型名称
  3. 获取用户组信息（Token 覆盖或用户默认）
  4. 支持 "auto" 组（按顺序尝试用户可访问的组）
  5. 调用 service.CacheGetRandomSatisfiedChannel() 查询渠道
  6. 按优先级和权重随机选择渠道
  7. 设置渠道上下文信息

## 5. 关键逻辑分析
- **能力表查询**: 通过 Ability 表 `(group, model, channel_id)` 实现 O(1) 渠道查找
- **优先级选择**: 渠道按优先级降序选择，同优先级内按权重随机
- **亲和性支持**: 记录用户与渠道的关联，后续请求优先使用相同渠道
- **自动组**: 支持 "auto" 组，按顺序尝试用户可访问的组

## 6. 关联文件
- `model/ability.go` - 能力表模型
- `service/channel.go` - 渠道缓存管理
- `service/channel_select.go` - 渠道选择算法
- `service/channel_affinity.go` - 渠道亲和性
