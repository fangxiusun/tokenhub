# payment_setting.go 代码阅读文档

## 1. 全局总结

该文件定义支付设置，包括充值金额选项、折扣配置、合规条款确认信息。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 常量名 | 值 | 说明 |
|--------|------|------|
| `CurrentComplianceTermsVersion` | `"v1"` | 当前合规条款版本 |

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `PaymentSetting` | `AmountOptions` | `[]int` | 充值金额选项 |
| | `AmountDiscount` | `map[int]float64` | 金额→折扣映射 |
| | `ComplianceConfirmed` | `bool` | 是否已确认合规条款 |
| | `ComplianceTermsVersion` | `string` | 合规条款版本 |
| | `ComplianceConfirmedAt` | `int64` | 确认时间戳 |
| | `ComplianceConfirmedBy` | `int` | 确认人 ID |
| | `ComplianceConfirmedIP` | `string` | 确认人 IP |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetPaymentSetting` | `func GetPaymentSetting() *PaymentSetting` | 获取支付设置 |
| `IsPaymentComplianceConfirmed` | `func IsPaymentComplianceConfirmed() bool` | 检查合规条款是否已确认 |

## 5. 关键逻辑分析

- 默认充值金额选项：10、20、50、100、200、500
- 合规确认需同时满足 `ComplianceConfirmed` 且版本匹配

## 6. 关联文件

- `setting/operation_setting/payment_setting_old.go` — 旧版支付设置
- `controller/payment.go` — 支付接口
