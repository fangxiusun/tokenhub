# settle.go 代码阅读文档

## 1. 全局总结
该文件实现了计费表达式的结算逻辑。它提供 ComputeTieredQuota 和 ComputeTieredQuotaWithRequest 函数，根据预消费时冻结的计费快照（BillingSnapshot）和实际 token 参数计算最终配额。包含配额转换、组比例应用和层级跨越检测。

## 2. 依赖关系
无外部依赖，仅使用本包内的函数和类型。

## 3. 类型定义
无新的类型定义，但使用了 billingexpr 包中的 BillingSnapshot、TokenParams、RequestInput、TieredResult、TraceResult 等类型。

## 4. 函数详解
### quotaConversion
内部函数，将表达式输出转换为配额。根据表达式版本（当前仅 v1）应用转换公式：exprOutput / 1,000,000 * QuotaPerUnit。

### ComputeTieredQuota
结算函数，接受 BillingSnapshot 和 TokenParams，返回 TieredResult。内部调用 ComputeTieredQuotaWithRequest。

### ComputeTieredQuotaWithRequest
类似 ComputeTieredQuota，但接受 RequestInput 参数，支持请求探测。执行表达式，计算配额，应用组比例，检测层级跨越。

## 5. 关键逻辑分析
- **结算流程**: 1) 使用预计算的哈希执行表达式；2) 将表达式输出转换为配额；3) 应用组比例并四舍五入；4) 比较实际层级与估计层级，检测是否跨越层级。
- **配额转换**: 公式基于表达式版本，当前 v1 假设表达式系数为 $/1M tokens 价格。
- **层级跨越检测**: 如果实际匹配的层级与预消费时估计的层级不同，则标记 CrossedTier 为 true，用于审计和调整。

## 6. 关联文件
- **run.go**: 提供 RunExprByHashWithRequest 函数，执行表达式。
- **types.go**: 定义 BillingSnapshot、TieredResult 等类型。
- **round.go**: 提供 QuotaRound 函数，用于配额四舍五入。
- **compile.go**: 提供 ExprHashString 函数，用于计算表达式哈希。
- **billingexpr_test.go**: 测试结算逻辑的各种场景。