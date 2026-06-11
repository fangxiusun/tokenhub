# checkin_setting.go 代码阅读文档

## 1. 全局总结

该文件定义签到功能的配置，包括开关、最小/最大额度奖励范围。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `CheckinSetting` | `Enabled` | `bool` | `false` | 是否启用签到功能 |
| | `MinQuota` | `int` | `1000` | 最小额度奖励（约 0.002 USD） |
| | `MaxQuota` | `int` | `10000` | 最大额度奖励（约 0.02 USD） |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetCheckinSetting` | `func GetCheckinSetting() *CheckinSetting` | 获取签到配置 |
| `IsCheckinEnabled` | `func IsCheckinEnabled() bool` | 检查签到功能是否启用 |
| `GetCheckinQuotaRange` | `func GetCheckinQuotaRange() (min, max int)` | 获取签到额度范围 |

## 5. 关键逻辑分析

- 签到功能默认关闭
- 额度范围：1000-10000（约 0.002-0.02 USD）

## 6. 关联文件

- `service/checkin.go` — 签到服务实现
- `controller/checkin.go` — 签到接口
