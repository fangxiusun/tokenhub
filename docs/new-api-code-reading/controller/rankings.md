# rankings.go 代码阅读文档

## 1. 全局总结

该文件提供用户/渠道/模型的排行榜查询接口。

## 2. 依赖关系

- `service` — 排行榜快照获取
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetRankings(c *gin.Context)`
获取排行榜数据。支持按时间周期过滤（默认 week）。

## 5. 关键逻辑分析

时间周期参数：day、week、month 等。

## 6. 关联文件

- `service/rankings.go` — 排行榜实现
