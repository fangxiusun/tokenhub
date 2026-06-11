# usedata.go 代码阅读文档

## 1. 全局总结

该文件提供了用户使用数据的统计查询接口，包括按日期的额度使用统计。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 使用数据模型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetAllQuotaDates(c *gin.Context)`
管理员获取所有用户的按日期额度使用统计。支持时间范围和用户名过滤。

### `GetQuotaDatesByUser(c *gin.Context)`
获取当前用户的按日期额度使用统计。

## 5. 关键逻辑分析

- 按日期聚合额度使用数据
- 支持时间范围过滤

## 6. 关联文件

- `model/usedata.go` — 使用数据模型
