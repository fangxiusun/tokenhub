# types.go 代码阅读文档

## 1. 全局总结
该文件定义了计费表达式系统的核心数据类型，包括请求输入、token 参数、追踪结果、计费快照和分层结果。同时提供表达式哈希函数。

## 2. 依赖关系
- **crypto/sha256**: 用于计算表达式字符串的 SHA-256 哈希。
- **fmt**: 格式化哈希值。

## 3. 类型定义
### RequestInput
结构体，表示请求输入：
- **Headers**: 请求头键值对。
- **Body**: 请求体字节数组。

### TokenParams
结构体，表示所有 token 维度：
- **P**: 提示 token 数（文本），自动排除单独定价的子类别。
- **C**: 完成 token 数（文本），自动排除单独定价的子类别。
- **Len**: 总输入上下文长度，用于层级条件（非 Claude：原始 prompt_tokens；Claude：文本 + 缓存读取 + 缓存创建）。
- **CR**: 缓存读取（命中）token 数。
- **CC**: 缓存创建 token 数（Claude 为 5 分钟 TTL，其他为通用）。
- **CC1h**: 缓存创建 token 数（仅 Claude，1 小时 TTL）。
- **Img**: 图像输入 token 数。
- **ImgO**: 图像输出 token 数。
- **AI**: 音频输入 token 数。
- **AO**: 音频输出 token 数。

### TraceResult
结构体，表示层级追踪结果：
- **MatchedTier**: 匹配的层级名称。
- **Cost**: 计算成本。

### BillingSnapshot
结构体，表示预消费时冻结的计费规则状态：
- **BillingMode**: 计费模式（如 "tiered_expr"）。
- **ModelName**: 模型名称。
- **ExprString**: 表达式字符串。
- **ExprHash**: 表达式哈希。
- **GroupRatio**: 组比例。
- **EstimatedPromptTokens**: 估计提示 token 数。
- **EstimatedCompletionTokens**: 估计完成 token 数。
- **EstimatedQuotaBeforeGroup**: 估计组前配额。
- **EstimatedQuotaAfterGroup**: 估计组后配额。
- **EstimatedTier**: 估计层级。
- **QuotaPerUnit**: 每单位配额。
- **ExprVersion**: 表达式版本。

### TieredResult
结构体，表示分层结算结果：
- **ActualQuotaBeforeGroup**: 实际组前配额。
- **ActualQuotaAfterGroup**: 实际组后配额（整数）。
- **MatchedTier**: 实际匹配的层级。
- **CrossedTier**: 是否跨越层级。

## 4. 函数详解
### ExprHashString
计算表达式字符串的 SHA-256 十六进制摘要。用于缓存键和快照存储。

## 5. 关键逻辑分析
- **类型设计**: 类型设计考虑了向后兼容性（可选字段默认为 0）、可序列化性（BillingSnapshot 无指针）和扩展性（预留版本字段）。
- **哈希确定性**: 使用 SHA-256 确保相同表达式产生相同哈希，用于缓存和快照匹配。

## 6. 关联文件
- **compile.go**: 使用 ExprHashString 计算缓存键。
- **run.go**: 使用 TokenParams 和 RequestInput 作为执行参数。
- **settle.go**: 使用 BillingSnapshot 和 TieredResult 进行结算。
- **billingexpr_test.go**: 测试这些类型的功能。