# ratio_config.go 代码阅读文档

## 1. 全局总结

该文件提供倍率配置的公开查询接口。需要启用倍率暴露开关才能访问。

## 2. 依赖关系

- `setting/ratio_setting` — 倍率配置和暴露控制
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `GetRatioConfig(c *gin.Context)`
获取倍率配置数据。需要 `IsExposeRatioEnabled()` 为 true。

## 5. 关键逻辑分析

- 通过开关控制是否对外暴露倍率配置
- 未启用时返回 403

## 6. 关联文件

- `setting/ratio_setting/` — 倍率配置
