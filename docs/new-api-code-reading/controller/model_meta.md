# model_meta.go 代码阅读文档

## 1. 全局总结

该文件实现了模型元数据（Model Meta）的 CRUD 管理接口。支持分页查询、搜索、创建、更新和删除模型元数据，并批量填充附加信息（端点、渠道、分组、计费类型）。

## 2. 依赖关系

- `common` — 通用工具、分页
- `constant` — 端点类型
- `model` — 模型元数据模型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetAllModelsMeta(c *gin.Context)`
分页获取所有模型元数据。返回列表、总数和供应商计数。

### `SearchModelsMeta(c *gin.Context)`
按关键词和供应商搜索模型。

### `GetModelMeta(c *gin.Context)`
根据 ID 获取单个模型详情。

### `CreateModelMeta(c *gin.Context)`
创建新模型。检查名称重复，创建后刷新定价缓存。

### `UpdateModelMeta(c *gin.Context)`
更新模型。支持 `status_only` 模式仅更新状态。更新后刷新定价缓存。

### `DeleteModelMeta(c *gin.Context)`
删除模型。删除后刷新定价缓存。

### `enrichModels(models []*model.Model)`
批量填充模型附加信息。区分精确匹配和规则匹配模型，使用批量查询避免 N+1 问题。

## 5. 关键逻辑分析

- 规则匹配模型（前缀/后缀/包含）需要遍历定价缓存进行内存匹配
- 批量查询渠道绑定信息（`GetBoundChannelsByModelsMap`）减少数据库查询
- 规则模型的端点/分组/配额类型取所有匹配模型的并集
- 创建/更新/删除后都会调用 `model.RefreshPricing()` 刷新缓存

## 6. 关联文件

- `model/model.go` — 模型元数据模型
- `model/pricing.go` — 定价缓存
