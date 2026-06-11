# dto.go 代码阅读文档

## 1. 全局总结
本文件定义了 SiliconFlow API 专用的数据传输对象，包括 Token 统计、Rerank 响应和图像请求的结构体。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/dto` | 引用 `dto.RerankResponseResult` 类型 |

## 3. 类型定义

### `SFTokens`
SiliconFlow Token 统计结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `InputTokens` | `int` | `json:"input_tokens"` | 输入 token 数量 |
| `OutputTokens` | `int` | `json:"output_tokens"` | 输出 token 数量 |

### `SFMeta`
SiliconFlow 元数据结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Tokens` | `SFTokens` | `json:"tokens"` | Token 统计信息 |

### `SFRerankResponse`
SiliconFlow Rerank 响应结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Results` | `[]dto.RerankResponseResult` | `json:"results"` | Rerank 结果列表 |
| `Meta` | `SFMeta` | `json:"meta"` | 元数据（含 token 统计） |

### `SFImageRequest`
SiliconFlow 图像请求结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Model` | `string` | `json:"model"` | 模型名称 |
| `Prompt` | `string` | `json:"prompt"` | 生成提示 |
| `NegativePrompt` | `string` | `json:"negative_prompt,omitempty"` | 反向提示 |
| `ImageSize` | `string` | `json:"image_size,omitempty"` | 图像尺寸 |
| `BatchSize` | `uint` | `json:"batch_size,omitempty"` | 批量生成数量 |
| `Seed` | `uint64` | `json:"seed,omitempty"` | 随机种子 |
| `NumInferenceSteps` | `uint` | `json:"num_inference_steps,omitempty"` | 推理步数 |
| `GuidanceScale` | `float64` | `json:"guidance_scale,omitempty"` | 引导比例 |
| `Cfg` | `float64` | `json:"cfg,omitempty"` | CFG 缩放因子 |
| `Image` | `string` | `json:"image,omitempty"` | 参考图像 URL |
| `Image2` | `string` | `json:"image2,omitempty"` | 第二参考图像 |
| `Image3` | `string` | `json:"image3,omitempty"` | 第三参考图像 |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- **Rerank Token 统计**：SiliconFlow 的 Rerank 响应包含 token 使用统计，格式与聊天 API 不同（使用 `meta.tokens` 而非 `usage`）。
- **图像参数丰富**：`SFImageRequest` 支持多种图像生成参数，包括反向提示、推理步数、引导比例等。
- **多图像输入**：支持最多三个参考图像输入（`Image`、`Image2`、`Image3`），用于图像编辑或风格迁移场景。

## 6. 关联文件
- `relay/channel/siliconflow/adaptor.go` — 使用这些结构体
- `relay/channel/siliconflow/relay-siliconflow.go` — Rerank 响应处理
- `dto/rerank.go` — `RerankResponseResult` 类型定义
