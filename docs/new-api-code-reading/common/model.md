# model.go 代码阅读文档

## 1. 全局总结

`model.go` 是模型分类工具文件，定义了 AI 模型的分类列表和判断函数。主要用于识别模型类型（如仅 OpenAI Responses API 模型、图像生成模型、文本模型），支持模型路由、功能过滤和 API 分发。

## 2. 依赖关系

### 标准库依赖
- `strings` - 字符串匹配和处理

### 项目内部依赖
- 无

## 3. 类型定义

### 全局变量（模型列表）

```go
var OpenAIResponseOnlyModels []string  // 仅支持 OpenAI Responses API 的模型
var ImageGenerationModels []string     // 图像生成模型
var OpenAITextModels []string          // OpenAI 文本模型
```

#### OpenAIResponseOnlyModels
仅可用于 OpenAI Responses API 的模型列表：
- `o3-pro`
- `o3-deep-research`
- `o4-mini-deep-research`

#### ImageGenerationModels
图像生成模型列表：
- `dall-e-3` - DALL-E 3
- `dall-e-2` - DALL-E 2
- `gpt-image-1` - GPT Image 1
- `prefix:imagen-` - Imagen 系列（前缀匹配）
- `flux-` - Flux 系列
- `flux.1-` - Flux 1.x 系列

#### OpenAITextModels
OpenAI 文本模型前缀：
- `gpt-` - GPT 系列
- `o1` - o1 推理模型
- `o3` - o3 推理模型
- `o4` - o4 推理模型
- `chatgpt` - ChatGPT 相关模型

## 4. 函数详解

### IsOpenAIResponseOnlyModel(modelName string) bool
```go
func IsOpenAIResponseOnlyModel(modelName string) bool
```
判断模型是否仅支持 OpenAI Responses API。

**参数：**
- `modelName` - 模型名称

**返回值：**
- `bool` - 是否为仅 Responses API 模型

**匹配方式：** 子字符串包含匹配

### IsImageGenerationModel(modelName string) bool
```go
func IsImageGenerationModel(modelName string) bool
```
判断模型是否为图像生成模型。

**参数：**
- `modelName` - 模型名称

**返回值：**
- `bool` - 是否为图像生成模型

**实现逻辑：**
1. 将模型名称转换为小写
2. 遍历图像生成模型列表
3. 支持两种匹配方式：
   - 子字符串包含匹配
   - `prefix:` 前缀匹配（忽略大小写）

### IsOpenAITextModel(modelName string) bool
```go
func IsOpenAITextModel(modelName string) bool
```
判断模型是否为 OpenAI 文本模型。

**参数：**
- `modelName` - 模型名称

**返回值：**
- `bool` - 是否为文本模型

**匹配方式：** 子字符串包含匹配（忽略大小写）

## 5. 关键逻辑分析

### 模型匹配策略

| 函数 | 匹配方式 | 大小写敏感 | 示例 |
|------|---------|-----------|------|
| `IsOpenAIResponseOnlyModel` | 包含 | 敏感 | `o3-pro` 匹配 `o3-pro` |
| `IsImageGenerationModel` | 包含/前缀 | 不敏感 | `flux-schnell` 匹配 `flux-` |
| `IsOpenAITextModel` | 包含 | 不敏感 | `GPT-4` 匹配 `gpt-` |

### 前缀匹配实现
```go
if strings.HasPrefix(m, "prefix:") && strings.HasPrefix(modelName, strings.TrimPrefix(m, "prefix:")) {
    return true
}
```
- 使用 `prefix:` 标记表示前缀匹配模式
- `prefix:imagen-` 表示匹配所有以 `imagen-` 开头的模型

### 大小写处理
```go
modelName = strings.ToLower(modelName)  // 转换为小写后匹配
```
- `IsImageGenerationModel` 和 `IsOpenAITextModel` 不区分大小写
- `IsOpenAIResponseOnlyModel` 保持原始大小写

## 6. 关联文件

- `relay/adaptor.go` - 模型适配器路由，使用模型分类决定处理逻辑
- `relay/controller.go` - 请求处理，使用模型分类进行功能过滤
- `model/channel.go` - 渠道模型，可能引用模型分类
- `setting/model.go` - 模型配置，管理模型列表
- `router/api.go` - API 路由，可能使用模型分类进行分发
