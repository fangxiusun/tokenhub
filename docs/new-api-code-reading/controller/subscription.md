# subscription.go 代码阅读文档

## 1. 全局总结

该文件实现了订阅计划（Subscription Plan）的完整管理，包括用户端的计划查看、余额购买、偏好设置，以及管理端的计划 CRUD、状态管理、用户绑定/解绑。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 订阅计划和用户订阅模型
- `setting/operation_setting` — 支付合规
- `setting/ratio_setting` — 组配置
- `gin-gonic/gin` — HTTP 框架
- `gorm.io/gorm` — 数据库事务

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `SubscriptionPlanDTO` | 计划数据传输对象 |
| `BillingPreferenceRequest` | 计费偏好请求 |
| `SubscriptionBalancePayRequest` | 余额购买请求 |
| `AdminUpsertSubscriptionPlanRequest` | 管理员创建/更新请求 |
| `AdminBindSubscriptionRequest` | 管理员绑定请求 |

## 4. 函数详解

### 用户 API
- `GetSubscriptionPlans` — 获取已启用的计划列表
- `GetSubscriptionSelf` — 获取当前用户的订阅信息（含已过期）
- `UpdateSubscriptionPreference` — 更新计费偏好（quota/subscription）
- `SubscriptionRequestBalancePay` — 使用余额购买订阅

### 管理员 API
- `AdminListSubscriptionPlans` — 获取所有计划
- `AdminCreateSubscriptionPlan` — 创建计划
- `AdminUpdateSubscriptionPlan` — 更新计划（事务操作）
- `AdminUpdateSubscriptionPlanStatus` — 启用/禁用计划
- `AdminBindSubscription` — 绑定用户订阅
- `AdminListUserSubscriptions` — 查看用户订阅
- `AdminCreateUserSubscription` — 为用户创建订阅
- `AdminInvalidateUserSubscription` — 取消用户订阅
- `AdminDeleteUserSubscription` — 删除用户订阅

## 5. 关键逻辑分析

- 计划价格上限 9999 USD
- 货币固定为 USD
- 支持多种时长单位：month、year、custom
- 预设折扣（AmountDiscount）按金额匹配
- 需要支付合规确认才能创建/更新计划

## 6. 关联文件

- `model/subscription.go` — 订阅模型
- `controller/payment_compliance.go` — 支付合规
