# channel_affinity_setting.go 代码阅读文档

## 1. 全局总结

该文件实现渠道亲和性（Channel Affinity）配置，允许根据请求特征（模型、路径、用户代理等）将请求绑定到特定渠道，支持 CLI 工具的请求头透传。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `ChannelAffinityKeySource` | `Type` | `string` | 键来源类型（context_int, context_string, request_header, gjson） |
| | `Key` | `string` | 键名 |
| | `Path` | `string` | gjson 路径 |
| `ChannelAffinityRule` | `Name` | `string` | 规则名称 |
| | `ModelRegex` | `[]string` | 模型匹配模式 |
| | `PathRegex` | `[]string` | 路径匹配模式 |
| | `UserAgentInclude` | `[]string` | 用户代理包含 |
| | `KeySources` | `[]ChannelAffinityKeySource` | 键来源列表 |
| | `ValueRegex` | `string` | 值匹配正则 |
| | `TTLSeconds` | `int` | 缓存 TTL |
| | `ParamOverrideTemplate` | `map[string]interface{}` | 参数覆盖模板 |
| | `SkipRetryOnFailure` | `bool` | 失败时跳过重试 |
| | `IncludeUsingGroup` | `bool` | 包含使用分组 |
| | `IncludeModelName` | `bool` | 包含模型名 |
| | `IncludeRuleName` | `bool` | 包含规则名 |
| `ChannelAffinitySetting` | `Enabled` | `bool` | 是否启用 |
| | `SwitchOnSuccess` | `bool` | 成功后切换 |
| | `KeepOnChannelDisabled` | `bool` | 渠道禁用时保持 |
| | `MaxEntries` | `int` | 最大缓存条目数 |
| | `DefaultTTLSeconds` | `int` | 默认 TTL |
| | `Rules` | `[]ChannelAffinityRule` | 规则列表 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `buildPassHeaderTemplate` | `func buildPassHeaderTemplate(headers []string) map[string]interface{}` | 构建请求头透传模板 |
| `GetChannelAffinitySetting` | `func GetChannelAffinitySetting() *ChannelAffinitySetting` | 获取渠道亲和性配置 |

## 5. 关键逻辑分析

- 默认规则支持 Codex CLI 和 Claude CLI 的请求追踪
- Codex CLI 透传头：`Originator`、`Session_id`、`User-Agent` 等
- Claude CLI 透传头：`X-Stainless-*`、`User-Agent`、`Anthropic-*` 等
- 使用 gjson 路径从请求体提取键值

## 6. 关联文件

- `relay/handler.go` — 使用渠道亲和性
- `middleware/distributor.go` — 渠道分发
