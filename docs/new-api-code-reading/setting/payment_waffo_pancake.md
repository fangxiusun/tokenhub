# payment_waffo_pancake.go 代码阅读文档

## 1. 全局总结

该文件定义 Waffo Pancake 托管结账的配置变量。当 MerchantID + PrivateKey + ProductID 非空时自动启用，无需单独的 Enabled 开关。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `WaffoPancakeMerchantID` | `string` | `""` | 商户 ID |
| `WaffoPancakePrivateKey` | `string` | `""` | 私钥 |
| `WaffoPancakeReturnURL` | `string` | `""` | 返回 URL |
| `WaffoPancakeUnitPrice` | `float64` | `1.0` | 单价 |
| `WaffoPancakeMinTopUp` | `int` | `1` | 最小充值金额 |
| `WaffoPancakeStoreID` | `string` | `""` | 商店 ID |
| `WaffoPancakeProductID` | `string` | `""` | 产品 ID |

## 4. 函数详解

无函数定义，仅变量声明。

## 5. 关键逻辑分析

- 网关启用逻辑：当 MerchantID + PrivateKey + ProductID 非空时自动启用
- StoreID 和 ProductID 由运营人员通过 `SaveWaffoPancakeConfig` 绑定

## 6. 关联文件

- `setting/payment_waffo.go` — Waffo 支付基础配置
- `service/waffo_pancake.go` — Waffo Pancake 支付服务
