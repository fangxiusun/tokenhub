# missing_models.go 代码阅读文档

## 1. 全局总结

该文件提供查询缺失模型列表的接口。缺失模型是指渠道中引用但模型元数据表中没有对应记录的模型名称。

## 2. 依赖关系

- `model` — 缺失模型查询
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetMissingModels(c *gin.Context)`
返回所有缺失模型名称列表。

## 5. 关键逻辑分析

辅助管理员快速发现需要配置的模型，通常配合 `model_sync.go` 使用。

## 6. 关联文件

- `controller/model_sync.go` — 模型同步
- `model/model.go` — `GetMissingModels` 查询
