# channel_cache.go 代码阅读文档

## 1. 全局总结
channel_cache.go 实现了渠道缓存管理，使用 Redis 和内存双重缓存提高渠道查询性能。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具（Redis）
- `constant` - 常量定义

### 2.2 被引用的文件
- `service/channel.go` - 渠道缓存管理
- `middleware/distributor.go` - 渠道路由分发

## 3. 类型定义
### 3.1 变量
- `channelCache` - 内存缓存
- `channelCacheMap` - 渠道映射缓存

## 4. 函数详解
### 4.1 CacheGetRandomSatisfiedChannel
- **职责**: 从缓存中获取随机满足条件的渠道

### 4.2 CacheGetChannelById
- **职责**: 根据 ID 获取渠道

### 4.3 CacheUpdateChannel
- **职责**: 更新渠道缓存

## 5. 关键逻辑分析
- **双重缓存**: Redis 优先，内存回退
- **批量加载**: 首次查询时批量加载所有渠道
- **缓存失效**: 渠道更新时清除相关缓存

## 6. 关联文件
- `common/redis.go` - Redis 工具
- `service/channel.go` - 渠道管理
