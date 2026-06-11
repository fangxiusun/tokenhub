# general_setting.go 代码阅读文档

## 1. 全局总结

该文件定义通用操作设置，包括文档链接、Ping 间隔、额度展示类型（USD/CNY/TOKENS/CUSTOM）和自定义货币汇率。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 常量名 | 值 | 说明 |
|--------|------|------|
| `QuotaDisplayTypeUSD` | `"USD"` | 美元展示 |
| `QuotaDisplayTypeCNY` | `"CNY"` | 人民币展示 |
| `QuotaDisplayTypeTokens` | `"TOKENS"` | Token 数量展示 |
| `QuotaDisplayTypeCustom` | `"CUSTOM"` | 自定义货币展示 |

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `GeneralSetting` | `DocsLink` | `string` | 文档链接 |
| | `PingIntervalEnabled` | `bool` | 是否启用 Ping 间隔 |
| | `PingIntervalSeconds` | `int` | Ping 间隔秒数 |
| | `QuotaDisplayType` | `string` | 额度展示类型 |
| | `CustomCurrencySymbol` | `string` | 自定义货币符号 |
| | `CustomCurrencyExchangeRate` | `float64` | 自定义货币与美元汇率 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetGeneralSetting` | `func GetGeneralSetting() *GeneralSetting` | 获取通用设置 |
| `IsCurrencyDisplay` | `func IsCurrencyDisplay() bool` | 是否以货币形式展示 |
| `IsCNYDisplay` | `func IsCNYDisplay() bool` | 是否以人民币展示 |
| `GetQuotaDisplayType` | `func GetQuotaDisplayType() string` | 获取额度展示类型 |
| `GetCurrencySymbol` | `func GetCurrencySymbol() string` | 获取当前货币符号 |
| `GetUsdToCurrencyRate` | `func GetUsdToCurrencyRate(usdToCny float64) float64` | 获取 USD 到目标货币的汇率 |

## 5. 关键逻辑分析

- `GetCurrencySymbol` 根据展示类型返回对应符号：`$`、`¥`、自定义符号、或空字符串
- `GetUsdToCurrencyRate` 对于 TOKENS 类型不适用，返回 1
- 自定义货币汇率默认 1.0

## 6. 关联文件

- `controller/option.go` — 管理界面配置接口
- `web/default/src/` — 前端额度展示
