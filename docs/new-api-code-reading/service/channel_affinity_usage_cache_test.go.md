# channel_affinity_usage_cache_test.go 代码阅读文档

## 1. 全局总结
channel_affinity_usage_cache_test.go 测试渠道亲和性使用量缓存功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `service/channel_affinity.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestUsageCache` - 测试使用量缓存

## 4. 函数详解
### 4.1 TestUsageCache
- **职责**: 测试使用量缓存的读写操作
- **验证点**: 缓存命中、过期、更新

## 5. 关键逻辑分析
- **缓存策略**: 测试缓存的过期和更新机制

## 6. 关联文件
- `service/channel_affinity.go` - 被测试的文件
