# pricing_refresh.go 代码阅读文档

## 1. 全局总结

`pricing_refresh.go` 提供一个强制立即刷新定价缓存的函数，供内部管理 API 使用。它绕过默认的 1 分钟延迟刷新机制，直接获取锁并调用核心刷新逻辑，确保管理操作能立即看到最新数据。

## 2. 依赖关系

该文件无外部导入，仅使用同包内的 `updatePricingLock`、`modelSupportEndpointsLock` 和 `updatePricing()` 函数。

## 3. 类型定义

无类型定义。

## 4. 函数详解

### RefreshPricing()
强制刷新定价缓存。依次获取 `updatePricingLock` 和 `modelSupportEndpointsLock` 两个互斥锁，然后调用 `updatePricing()` 执行完整的缓存重建。与 `GetPricing()` 不同，该函数不检查时间间隔，直接执行刷新。

## 5. 关键逻辑分析

**锁顺序一致性**：`RefreshPricing()` 按 `updatePricingLock` → `modelSupportEndpointsLock` 的顺序获取锁，与 `GetPricing()` 中的锁顺序一致，避免死锁。

**适用场景**：管理员修改定价配置后调用此函数，确保变更立即生效。正常查询使用 `GetPricing()` 的延迟刷新即可。

## 6. 关联文件

- `pricing.go`：提供 `updatePricingLock`、`modelSupportEndpointsLock` 和 `updatePricing()`
