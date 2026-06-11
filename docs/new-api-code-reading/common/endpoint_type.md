# endpoint_type.go 代码阅读文档

## 1. 全局总结

本文件实现了根据渠道类型和模型名称确定 API 端点类型的功能。核心函数 `GetEndpointTypesByChannelType` 根据不同的 AI 渠道提供商返回其支持的端点类型列表，并按优先级排序。支持图像生成模型的端点类型会被插入到列表最前面。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/constant` | 渠道类型和端点类型常量 |

**被依赖方**：被端点路由选择和请求分发代码调用。

## 3. 类型定义

本文件无自定义类型定义，仅使用 `constant` 包中的枚举类型。

## 4. 函数详解

### `GetEndpointTypesByChannelType(channelType int, modelName string) []constant.EndpointType`

**功能**：根据渠道类型和模型名称获取支持的端点类型列表（按优先级排序）。

**参数**：
- `channelType int` — 渠道类型（对应 `constant.ChannelType*` 常量）
- `modelName string` — 模型名称

**返回值**：`[]constant.EndpointType` — 支持的端点类型列表

**逻辑**：

#### 渠道类型匹配（switch-case）

| 渠道类型 | 返回的端点类型 | 优先级 |
|----------|----------------|--------|
| `ChannelTypeJina` | `JinaRerank` | 1 |
| `ChannelTypeAws` / `ChannelTypeAnthropic` | `Anthropic`, `OpenAI` | 1, 2 |
| `ChannelTypeVertexAi` / `ChannelTypeGemini` | `Gemini`, `OpenAI` | 1, 2 |
| `ChannelTypeOpenRouter` | `OpenAI` | 1 |
| `ChannelTypeXai` | `OpenAI`, `OpenAIResponse` | 1, 2 |
| `ChannelTypeSora` | `OpenAIVideo` | 1 |
| 默认（其他渠道） | `OpenAI` 或 `OpenAIResponse` | 1 |

#### 默认情况的特殊处理

对于未明确列出的渠道类型：
1. 调用 `IsOpenAIResponseOnlyModel(modelName)` 判断模型是否仅支持 Response API
   - 是 → 返回 `OpenAIResponse`
   - 否 → 返回 `OpenAI`（所有渠道都支持 OpenAI 端点）

#### 图像生成模型特殊处理

如果模型是图像生成模型（`IsImageGenerationModel(modelName)` 为 `true`），将 `ImageGeneration` 端点类型插入到列表最前面，使其优先级最高。

## 5. 关键逻辑分析

### 端点优先级设计

端点类型的顺序代表优先级，系统会按顺序尝试匹配：

```
1. 首选端点（如 Anthropic、Gemini）
2. 备选端点（如 OpenAI）
```

这种设计允许：
- 渠道原生 API 优先使用
- 在原生 API 不可用时自动降级到 OpenAI 兼容接口

### 图像生成模型的特殊处理

图像生成模型的端点类型被插入到列表最前面：

```go
endpointTypes = append([]constant.EndpointType{constant.EndpointTypeImageGeneration}, endpointTypes...)
```

这意味着对于图像生成模型，系统会优先尝试图像生成端点。

### 注释掉的渠道

代码中注释了以下渠道的端点类型定义：
- Midjourney / MidjourneyPlus
- SunoAPI
- Kling
- Jimeng

这些可能是已弃用或尚未实现的渠道。

### 模型名称感知

函数不仅考虑渠道类型，还考虑模型名称：
- `IsOpenAIResponseOnlyModel()` — 判断模型是否仅支持 OpenAI Response API
- `IsImageGenerationModel()` — 判断模型是否为图像生成模型

这使得同一个渠道可以为不同模型提供不同的端点类型。

## 6. 关联文件

| 文件 | 关联关系 |
|------|----------|
| `constant/channel.go` | 渠道类型常量定义 |
| `constant/endpoint.go` | 端点类型常量定义 |
| `endpoint_defaults.go` | 端点默认配置 |
| `relay/` | 根据端点类型选择对应的请求转换器 |
| `common/model_utils.go` | `IsOpenAIResponseOnlyModel()`, `IsImageGenerationModel()` |
