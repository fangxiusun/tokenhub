# payment_waffo.go 代码阅读文档

## 1. 全局总结

该文件定义 Waffo 支付网关的配置变量，并提供支付方式的读写接口。支持沙箱/生产环境切换。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/common` — JSON 操作、OptionMap 读写
- `github.com/QuantumNous/new-api/constant` — Waffo 支付方式常量

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `WaffoEnabled` | `bool` | 是否启用 Waffo 支付 |
| `WaffoApiKey` | `string` | 生产环境 API 密钥 |
| `WaffoPrivateKey` | `string` | 生产环境私钥 |
| `WaffoPublicCert` | `string` | 生产环境公钥证书 |
| `WaffoSandboxPublicCert` | `string` | 沙箱环境公钥证书 |
| `WaffoSandboxApiKey` | `string` | 沙箱环境 API 密钥 |
| `WaffoSandboxPrivateKey` | `string` | 沙箱环境私钥 |
| `WaffoSandbox` | `bool` | 是否使用沙箱环境 |
| `WaffoMerchantId` | `string` | 商户 ID |
| `WaffoNotifyUrl` | `string` | 异步通知 URL |
| `WaffoReturnUrl` | `string` | 同步返回 URL |
| `WaffoSubscriptionReturnUrl` | `string` | 订阅返回 URL |
| `WaffoCurrency` | `string` | 货币类型 |
| `WaffoUnitPrice` | `float64` | 单价，默认 1.0 |
| `WaffoMinTopUp` | `int` | 最小充值金额，默认 1 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetWaffoPayMethods` | `func GetWaffoPayMethods() []constant.WaffoPayMethod` | 从 OptionMap 读取 Waffo 支付方式配置 |
| `SetWaffoPayMethods` | `func SetWaffoPayMethods(methods []constant.WaffoPayMethod) error` | 序列化并更新 OptionMap 中的支付方式 |
| `copyDefaultWaffoPayMethods` | `func copyDefaultWaffoPayMethods() []constant.WaffoPayMethod` | 返回默认支付方式的深拷贝 |
| `WaffoPayMethods2JsonString` | `func WaffoPayMethods2JsonString() string` | 将默认支付方式序列化为 JSON 字符串 |

## 5. 关键逻辑分析

- `GetWaffoPayMethods` 使用 `common.OptionMapRWMutex` 保护的 OptionMap 读取配置
- 当 OptionMap 中无配置时，返回默认支付方式的深拷贝
- 使用 `common.UnmarshalJsonStr` 反序列化 JSON 字符串

## 6. 关联文件

- `constant/waffo.go` — Waffo 支付方式常量定义
- `common/option.go` — OptionMap 全局配置存储
- `service/waffo.go` — Waffo 支付服务实现
