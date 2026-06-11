# payment_method_guard_test.go 代码阅读文档

## 1. 全局总结
该文件是支付方式守卫（Payment Method Guard）功能的单元测试文件。测试验证了在支付处理过程中，当支付提供者与记录的支付方式不匹配时，系统是否能正确拒绝操作并保持数据一致性。测试覆盖了充值、订单状态更新、订阅订单完成和过期等多个场景。

## 2. 依赖关系
- **testing**: Go 标准测试包。
- **time**: 用于生成测试时间戳。
- **github.com/QuantumNous/new-api/common**: 提供状态常量（如 UserStatusEnabled、TopUpStatusPending）。
- **github.com/stretchr/testify/assert**: 测试断言库。
- **github.com/stretchr/testify/require**: 测试断言库（失败时立即停止）。

## 3. 类型定义
该文件没有定义新的类型或结构体。

## 4. 函数详解
### 辅助函数
- **insertUserForPaymentGuardTest(t *testing.T, id int, quota int)**
  - 插入测试用户。

- **insertSubscriptionPlanForPaymentGuardTest(t *testing.T, id int) *SubscriptionPlan**
  - 插入测试订阅计划。

- **insertSubscriptionOrderForPaymentGuardTest(t *testing.T, tradeNo string, userID int, planID int, paymentProvider string)**
  - 插入测试订阅订单。

- **insertTopUpForPaymentGuardTest(t *testing.T, tradeNo string, userID int, paymentProvider string)**
  - 插入测试充值记录。

- **getTopUpStatusForPaymentGuardTest(t *testing.T, tradeNo string) string**
  - 获取充值记录状态。

- **countUserSubscriptionsForPaymentGuardTest(t *testing.T, userID int) int64**
  - 统计用户订阅数量。

- **getUserQuotaForPaymentGuardTest(t *testing.T, userID int) int**
  - 获取用户配额。

### 测试函数
- **TestRechargeWaffoPancake_RejectsMismatchedPaymentMethod(t *testing.T)**
  - 测试 RechargeWaffoPancake 函数：当充值记录的支付方式与预期不匹配时，应拒绝充值。
  - 验证充值状态保持 Pending，用户配额不变。

- **TestUpdatePendingTopUpStatus_RejectsMismatchedPaymentProvider(t *testing.T)**
  - 测试 UpdatePendingTopUpStatus 函数：当支付提供者与记录不匹配时，应拒绝状态更新。
  - 测试用例：
    1. Stripe 过期：记录为 Creem，预期为 Stripe，目标状态为 Expired。
    2. Waffo 失败：记录为 Stripe，预期为 Waffo，目标状态为 Failed。
  - 验证错误类型为 ErrPaymentMethodMismatch，状态保持 Pending。

- **TestCompleteSubscriptionOrder_RejectsMismatchedPaymentProvider(t *testing.T)**
  - 测试 CompleteSubscriptionOrder 函数：当支付提供者与记录不匹配时，应拒绝完成订阅订单。
  - 验证订单状态保持 Pending，用户无新订阅，无充值记录。

- **TestExpireSubscriptionOrder_RejectsMismatchedPaymentProvider(t *testing.T)**
  - 测试 ExpireSubscriptionOrder 函数：当支付提供者与记录不匹配时，应拒绝过期订阅订单。
  - 验证订单状态保持 Pending。

## 5. 关键逻辑分析
- **支付方式一致性检查**: 测试验证了系统在支付处理过程中会检查支付提供者与记录是否一致，防止支付混淆或攻击。
- **数据完整性保护**: 当检测到不匹配时，系统拒绝操作并保持数据不变，确保数据一致性。
- **错误类型验证**: 使用 `ErrorIs` 验证返回的错误类型是否为 `ErrPaymentMethodMismatch`。
- **状态不变性验证**: 每个测试用例都验证了相关记录的状态在操作失败后保持不变。
- **事务隔离**: 测试使用 `truncateTables` 清空表，确保测试之间互不影响。

## 6. 关联文件
- **model/topup.go**: 包含充值相关的数据模型和函数（如 RechargeWaffoPancake、UpdatePendingTopUpStatus）。
- **model/subscription.go**: 包含订阅相关的数据模型和函数（如 CompleteSubscriptionOrder、ExpireSubscriptionOrder）。
- **model/user.go**: 包含用户数据模型（User 结构体）。
- **common/constants.go**: 包含状态常量定义。