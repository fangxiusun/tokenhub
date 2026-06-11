# waffo_pay_method.go 代码阅读文档

## 1. 全局概述

本文件定义了 Waffo 支付方式的结构体和默认支付方式列表。Waffo 是一个支付网关集成，本文件为其提供支付方式的前端展示和 API 参数映射配置。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### WaffoPayMethod 结构体

```go
type WaffoPayMethod struct {
    Name          string `json:"name"`            // 前端显示名称
    Icon          string `json:"icon"`            // 前端图标标识符
    PayMethodType string `json:"payMethodType"`  // Waffo API PayMethodType，可逗号分隔
    PayMethodName string `json:"payMethodName"`  // Waffo API PayMethodName，空值表示自动选择
}
```

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 默认支付方式列表

`DefaultWaffoPayMethods` 定义了三种默认支付方式：

| 名称 | 图标 | 类型 | API 名称 | 说明 |
|------|------|------|----------|------|
| Card | `/pay-card.png` | `CREDITCARD,DEBITCARD` | (空) | 信用卡/借记卡，Waffo 自动选择 |
| Apple Pay | `/pay-apple.png` | `APPLEPAY` | `APPLEPAY` | Apple Pay |
| Google Pay | `/pay-google.png` | `GOOGLEPAY` | `GOOGLEPAY` | Google Pay |

### 设计说明

- `PayMethodType` 支持逗号分隔的多类型值（如 `CREDITCARD,DEBITCARD`）
- `PayMethodName` 为空时，Waffo 会自动选择最合适的支付方式
- 图标路径使用前端静态资源路径

## 6. 相关文件

- `controller/payment.go` — 支付控制器
- `setting/operation_setting.go` — 运营设置中可能引用支付方式
- `web/default/src/` — 前端使用这些支付方式配置
