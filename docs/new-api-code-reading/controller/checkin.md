# checkin.go 代码阅读文档

## 1. 全局总结

该文件实现了用户每日签到功能，包括查看签到状态/历史和执行签到。签到可获得随机额度奖励。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志（额度格式化）
- `model` — 用户签到模型
- `setting/operation_setting` — 签到配置（启用状态、最小/最大额度）
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetCheckinStatus(c *gin.Context)`
获取用户签到状态和历史记录。支持按月份查询（默认当前月），返回签到配置和统计数据。

### `DoCheckin(c *gin.Context)`
执行用户签到。调用 `model.UserCheckin` 完成签到并获取奖励额度，记录系统日志。

## 5. 关键逻辑分析

- 签到功能可通过 `operation_setting.GetCheckinSetting().Enabled` 开关控制
- 签到奖励额度在 `MinQuota` 和 `MaxQuota` 之间随机
- 签到记录会写入系统日志

## 6. 关联文件

- `model/checkin.go` — 签到数据模型
- `setting/operation_setting/` — 签到配置
