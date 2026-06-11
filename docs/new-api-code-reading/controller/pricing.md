# pricing.go 代码阅读文档

## 1. 全局总结

该文件实现了定价信息查询接口和模型倍率重置功能。返回用户可用的定价列表、供应商信息、组倍率等。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 定价模型
- `service` — 用户可用组查询
- `setting/ratio_setting` — 倍率设置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetPricing(c *gin.Context)`
获取定价信息。返回：定价列表、供应商、组倍率、用户可用组、支持的端点、自动组、定价版本。

### `ResetModelRatio(c *gin.Context)`
重置模型倍率为默认值。

### `filterPricingByUsableGroups(pricing, usableGroup)`
按用户可用组过滤定价列表。支持 "all" 组匹配所有。

## 5. 关键逻辑分析

- 定价列表按用户可用组过滤
- 组倍率支持组间倍率覆盖（`GetGroupGroupRatio`）
- 定价版本使用固定哈希值

## 6. 关联文件

- `model/pricing.go` — 定价模型
- `setting/ratio_setting/` — 倍率配置
