# billing.go 代码阅读文档

## 1. 全局总结

本文件定义了 `BillingSettler` 接口，抽象了计费会话的生命周期操作，包括结算、退款、预扣额度查询和补充。由 `service.BillingSession` 实现，存储在 `RelayInfo` 上以避免循环引用。

## 2. 依赖关系

- `gin`: Context 类型

## 3. 类型定义

### `BillingSettler` 接口
```go
type BillingSettler interface {
    Settle(actualQuota int) error      // 结算，计算 delta 并调整资金来源
    Refund(c *gin.Context)             // 退还预扣费（异步）
    NeedsRefund() bool                 // 是否需要退款
    GetPreConsumedQuota() int          // 获取预扣额度
    Reserve(targetQuota int) error     // 补充预扣额度到目标值
}
```

## 4. 函数详解

本文件仅定义接口，无函数实现。

## 5. 关键逻辑分析

1. **幂等安全**: Refund 操作是幂等的，重复调用不会产生副作用
2. **异步退款**: 通过 gopool 异步执行退款操作
3. **补充机制**: Reserve 用于在任务状态变化时补充预扣额度

## 6. 关联文件

- `relay/common/relay_info.go`: RelayInfo 中的 Billing 字段
- `service/billing.go`: BillingSettler 的具体实现
