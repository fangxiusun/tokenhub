# prefill_group.go 代码阅读文档

## 1. 全局总结

该文件实现了预填组（Prefill Group）的 CRUD 管理接口。预填组用于前端表单的预填充选项。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 预填组模型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetPrefillGroups(c *gin.Context)`
获取预填组列表，支持按类型过滤。

### `CreatePrefillGroup(c *gin.Context)`
创建新预填组。检查名称和类型非空，检查名称重复。

### `UpdatePrefillGroup(c *gin.Context)`
更新预填组。检查 ID 和名称重复。

### `DeletePrefillGroup(c *gin.Context)`
删除预填组。

## 5. 关键逻辑分析

- 名称和类型不能为空
- 创建/更新时检查名称唯一性

## 6. 关联文件

- `model/prefill_group.go` — 预填组模型
