# tiered_settle.go 代码阅读文档

## 1. 全局总结

该文件实现分级计费（Tiered Billing）的结算逻辑，支持基于表达式的动态定价。将 usage 数据转换为分级计费参数，并计算实际额度。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `dto` | Usage 结构体 |
| `billingexpr` | 分级计费表达式引擎 |
| `relaycommon` | RelayInfo |

## 3. 类型定义

### `TieredResultWrapper`
`billingexpr.TieredResult` 的类型别名

## 4. 函数详解

### `BuildTieredTokenParams(usage, isClaudeUsageSemantic, usedVars) billingexpr.TokenParams`
构建分级计费参数：
- GPT 格式：prompt_tokens 包含所有子类别
- Claude 格式：input_tokens 仅文本
- 根据表达式使用的变量，减去相应子类别

### `TryTieredSettle(relayInfo, params) (ok, quota, result)`
尝试分级计费结算：
1. 检查是否使用 tiered_expr 模式
2. 调用 `billingexpr.ComputeTieredQuotaWithRequest` 计算
3. 失败时回退到预扣额度

## 5. 关键逻辑分析

1. **Token 归一化**：将 GPT/Claude 不同的 usage 格式归一化
2. **变量减除**：根据表达式引用的变量，从 P/C 中减去子类别
3. **输入长度**：len 用于分级条件评估（包含缓存 token）
4. **错误回退**：表达式计算失败时使用预扣额度

## 6. 关联文件

- `pkg/billingexpr/` — 分级计费表达式引擎
- `text_quota.go` — 文本额度计算
