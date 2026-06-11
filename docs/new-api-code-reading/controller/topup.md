# topup.go 代码阅读文档

## 1. 全局总结

该文件实现了充值（TopUp）功能，包括获取充值信息、易支付（Epay）下单和回调、金额查询、充值记录查询和管理员补单。支持多种支付渠道（Epay、Stripe、Creem、Waffo 等）。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 充值订单模型
- `service` — 回调地址
- `setting` — 支付渠道配置
- `setting/operation_setting` — 支付设置
- `Calcium-Ion/go-epay` — 易支付 SDK
- `gin-gonic/gin` — HTTP 框架
- `samber/lo` — 集合操作
- `shopspring/decimal` — 精确计算

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `EpayRequest` | 易支付请求 |
| `AmountRequest` | 金额查询请求 |
| `AdminCompleteTopupRequest` | 管理员补单请求 |
| `refCountedMutex` | 带引用计数的互斥锁 |

## 4. 函数详解

### 充值信息
- `GetTopUpInfo` — 获取充值配置（支付方式、金额选项、折扣等）

### Epay 支付
- `RequestEpay` — 创建 Epay 充值订单
- `EpayNotify` — Epay 支付回调处理
- `RequestAmount` — 查询指定金额的实际支付价格

### 记录查询
- `GetUserTopUps` — 用户充值记录
- `GetAllTopUps` — 管理员充值记录

### 管理员操作
- `AdminCompleteTopUp` — 管理员补单

## 5. 关键逻辑分析

- 支付金额计算：amount × price × topupGroupRatio × discount
- 订单号格式：`USR{userId}NO{random6}{timestamp}`
- 回调使用引用计数互斥锁防止并发补单
- 充值成功后额度增加：amount × QuotaPerUnit
- 支持 TOKENS 展示类型的金额转换
- 支付合规确认前不显示支付方式

## 6. 关联文件

- `model/topup.go` — 充值订单模型
- `controller/payment_webhook_availability.go` — 支付渠道可用性
- `controller/return_path.go` — 返回路径
