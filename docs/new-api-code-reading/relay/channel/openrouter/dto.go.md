# dto.go 代码阅读文档

## 1. 全局总结
本文件定义了 OpenRouter 渠道专用的数据传输对象（DTO），包括推理（Reasoning）请求结构和企业版响应结构。用于与 OpenRouter API 进行数据交互时的序列化和反序列化。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `encoding/json` | 用于 `json.RawMessage` 类型，处理原始 JSON 数据 |

## 3. 类型定义

### `RequestReasoning`
OpenRouter 推理请求配置结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Enabled` | `bool` | `json:"enabled"` | 是否启用推理功能 |
| `Effort` | `string` | `json:"effort,omitempty"` | 推理努力程度，可选 "high"、"medium"、"low"（OpenAI 风格） |
| `MaxTokens` | `int` | `json:"max_tokens,omitempty"` | 推理 token 上限（Anthropic 风格） |
| `Exclude` | `bool` | `json:"exclude,omitempty"` | 是否从响应中排除推理 token，默认 false |

**注意**：`Effort` 和 `MaxTokens` 互斥，只能设置其中一个。

### `OpenRouterEnterpriseResponse`
OpenRouter 企业版 API 的通用响应结构。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Data` | `json.RawMessage` | `json:"data"` | 响应数据，保持原始 JSON 格式以便灵活解析 |
| `Success` | `bool` | `json:"success"` | 请求是否成功 |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- `RequestReasoning` 使用 `omitempty` 标签确保未设置的字段不会出现在 JSON 中。
- `OpenRouterEnterpriseResponse.Data` 使用 `json.RawMessage` 类型，允许调用方按需解析不同的响应结构，增强了灵活性。
- 该文件仅定义数据结构，实际的请求/响应处理逻辑在其他文件中实现。

## 6. 关联文件
- `relay/channel/openrouter/constant.go` — 渠道常量定义
- `relay/channel/adapter.go` — `Adaptor` 接口定义
- `relay/channel/openai/` — OpenAI 渠道适配器，OpenRouter 兼容 OpenAI 格式
