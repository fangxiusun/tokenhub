# text_quota.go 代码阅读文档

## 1. 全局总结

该文件实现文本请求的额度计算和消费日志记录，是文本类 API 计费的核心实现。支持 OpenAI 和 Claude 两种 usage 语义，处理缓存 token、图片 token、音频 token、Web Search 调用等多种计费维度。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 额度单位 |
| `constant` | 上下文键名 |
| `dto` | Usage 结构体 |
| `model` | 额度操作 |
| `billingexpr` | 分级计费 |
| `relaycommon` | RelayInfo |
| `operation_setting` | 工具价格 |
| `decimal` | 精确计算 |
| `gopool` | 异步任务 |

## 3. 类型定义

### `textQuotaSummary`
额度计算摘要：所有 token 计数、倍率、价格、最终额度

## 4. 核心函数

### `calculateTextQuotaSummary(ctx, relayInfo, usage) textQuotaSummary`
文本额度计算：
1. 初始化倍率信息
2. 处理 OpenRouter Claude 缓存 token
3. 计算基础 token（减去缓存/图片/音频）
4. 应用各类倍率
5. 计算工具调用附加费
6. 应用其他倍率（视频折扣等）

### `calculateTextToolCallSurcharge(ctx, relayInfo, summary) decimal.Decimal`
工具调用附加费计算：
- Web Search 调用
- Claude Web Search 调用
- File Search 调用
- Image Generation 调用

### `PostTextConsumeQuota(ctx, relayInfo, usage, extraContent)`
文本请求后结算：
1. 记录使用量缓存统计
2. 计算额度摘要
3. 尝试分级计费
4. 更新用户/通道额度
5. 结算计费
6. 记录消费日志

### `composeTieredTextQuota(relayInfo, summary, tieredQuota, tieredResult) int`
组合分级计费额度（含工具附加费）

## 5. 关键逻辑分析

1. **双语义支持**：OpenAI（prompt_tokens 包含缓存）和 Claude（缓存单独计算）
2. **OpenRouter Claude**：特殊的缓存 token 处理
3. **工具附加费**：Web Search、File Search、Image Generation 独立计费
4. **缓存写入 token**：5m/1h 分离的缓存创建 token

## 6. 关联文件

- `quota.go` — 其他类型额度计算
- `tiered_settle.go` — 分级计费
- `log_info_generate.go` — 日志信息生成
