# channel_affinity_cache.go 代码阅读文档

## 1. 全局总结

该文件提供渠道亲和性缓存的管理接口，包括查看缓存统计、清除缓存和查看使用量缓存统计。渠道亲和性缓存用于记录用户/Key 与渠道之间的绑定关系，提升请求路由效率。

## 2. 依赖关系

- `service` — 缓存统计和清除操作
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetChannelAffinityCacheStats(c *gin.Context)`
获取渠道亲和性缓存的全局统计信息。

### `ClearChannelAffinityCache(c *gin.Context)`
清除渠道亲和性缓存。支持两种模式：
- `all=true`：清除所有缓存
- `rule_name=xxx`：清除指定规则的缓存

### `GetChannelAffinityUsageCacheStats(c *gin.Context)`
获取指定规则、用户组和 Key 指纹的使用量缓存统计。需要 `rule_name` 和 `key_fp` 参数。

## 5. 关键逻辑分析

- 所有接口都是只读或清除操作，不涉及业务状态变更
- 清除操作支持按规则名精确清除或全量清除

## 6. 关联文件

- `service/channel_affinity.go` — 缓存管理的实际实现
