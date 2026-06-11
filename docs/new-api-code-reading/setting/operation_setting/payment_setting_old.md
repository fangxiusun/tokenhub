# payment_setting_old.go 代码阅读文档

## 1. 全局总结

该文件是旧版支付设置文件，包含支付地址、Epay 配置、价格、支付方式等变量。新参数应添加到 `payment_setting.go` 中。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/common` — JSON 操作

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `PayAddress` | `string` | `""` | 支付地址 |
| `CustomCallbackAddress` | `string` | `""` | 自定义回调地址 |
| `EpayId` | `string` | `""` | Epay ID |
| `EpayKey` | `string` | `""` | Epay 密钥 |
| `Price` | `float64` | `7.3` | 价格（人民币/美元） |
| `MinTopUp` | `int` | `1` | 最小充值金额 |
| `USDExchangeRate` | `float64` | `7.3` | 美元汇率 |
| `PayMethods` | `[]map[string]string` | 支付方式列表 | 默认包含支付宝、微信、自定义1 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `UpdatePayMethodsByJsonString` | `func UpdatePayMethodsByJsonString(jsonString string) error` | 从 JSON 更新支付方式 |
| `PayMethods2JsonString` | `func PayMethods2JsonString() string` | 将支付方式序列化为 JSON |
| `ContainsPayMethod` | `func ContainsPayMethod(method string) bool` | 检查是否包含指定支付方式 |

## 5. 关键逻辑分析

- 该文件保留向后兼容，新配置应添加到 `payment_setting.go`
- 支付方式使用 `[]map[string]string` 结构

## 6. 关联文件

- `setting/operation_setting/payment_setting.go` — 新版支付设置
- `controller/payment.go` — 支付接口
