# price_test.go 代码阅读文档

## 1. 全局总结

本文件是 `price.go` 的单元测试，测试了分层计费表达式的价格计算逻辑。

## 2. 依赖关系

- `gin`: 测试上下文
- `billingexpr`: 分层计费表达式
- `config`: 全局配置

## 3. 测试用例详解

### `TestModelPriceHelperTieredUsesPreloadedRequestInput`
- 测试分层计费使用预加载的 RequestInput
- 验证 stream 参数影响 tier 选择
- 验证 QuotaToPreConsume 计算正确

## 4. 关键逻辑分析

1. **配置注入**: 通过 config 设置 billing_mode 和 billing_expr
2. **Tier 选择**: stream=true 选择 "stream" tier，p * 3 = 3000，除以 1M 再乘以 QuotaPerUnit = 1500
3. **快照验证**: TieredBillingSnapshot 包含完整的计费快照信息

## 5. 关联文件

- `relay/helper/price.go`: 被测试的源文件
- `pkg/billingexpr/`: 分层计费表达式引擎
