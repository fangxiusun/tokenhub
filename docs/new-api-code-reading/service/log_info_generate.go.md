# log_info_generate.go 代码阅读文档

## 1. 全局总结

该文件生成请求日志的 `other` 信息字段，记录计费相关的详细信息，包括模型倍率、分组倍率、缓存 token、订阅信息等。是日志系统的核心数据组装层。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 上下文操作 |
| `constant` | 上下文键名 |
| `dto` | Usage 结构体 |
| `billingexpr` | 分级计费信息 |
| `relaycommon` | RelayInfo |

## 3. 函数详解

### `GenerateTextOtherInfo(ctx, relayInfo, ...) map[string]interface{}`
生成文本请求的 other 信息：
- 模型/分组/补全倍率
- 缓存 token 信息
- 首次响应时间
- 管理员信息（通道、多密钥等）
- 参数覆盖审计
- 流状态
- 计费来源信息

### `GenerateClaudeOtherInfo(...)` / `GenerateAudioOtherInfo(...)` / `GenerateWssOtherInfo(...)`
特定格式的 other 信息生成

### `GenerateMjOtherInfo(...)` / `InjectTieredBillingInfo(...)`
Midjourney 和分级计费信息

### 内部辅助函数
- `appendBillingInfo` — 计费来源和订阅详情
- `appendRequestConversionChain` — 请求格式转换链
- `appendStreamStatus` — 流状态信息
- `appendParamOverrideInfo` — 参数覆盖审计

## 4. 关键逻辑分析

1. **分级计费**：支持 tiered_expr 模式的计费信息注入
2. **订阅详情**：记录预扣量、后调量、剩余额度等
3. **流状态**：记录流式响应的结束原因和错误信息

## 5. 关联文件

- `quota.go` / `text_quota.go` — 计费计算
- `model/log.go` — 日志记录
