# valid_request.go 代码阅读文档

## 1. 全局总结

本文件实现了 Relay 模块的请求验证逻辑，为每种请求格式提供对应的验证函数。是请求处理的第一道防线，确保请求参数合法。

## 2. 依赖关系

- `common`: JSON 反序列化
- `dto`: 请求 DTO
- `types`: 错误类型
- `relay/constant`: RelayMode

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `GetAndValidateRequest(c, format) (dto.Request, error)`
- 统一入口，根据 RelayFormat 分发到对应的验证函数

### 各格式验证函数

#### `GetAndValidateTextRequest(c, relayMode) (*dto.GeneralOpenAIRequest, error)`
- 验证 OpenAI 文本请求
- 支持 completions/chat_completions/embeddings/moderations/edits 模式
- 验证 max_tokens、model、messages 等字段

#### `GetAndValidateClaudeRequest(c) (*dto.ClaudeRequest, error)`
- 验证 Claude 请求，确保 messages 和 model 存在

#### `GetAndValidateGeminiRequest(c) (*dto.GeminiChatRequest, error)`
- 验证 Gemini 请求，确保 contents 存在

#### `GetAndValidateResponsesRequest(c) (*dto.OpenAIResponsesRequest, error)`
- 验证 Responses 请求，确保 model 和 input 存在

#### `GetAndValidOpenAIImageRequest(c, relayMode) (*dto.ImageRequest, error)`
- 验证图像请求
- 支持 multipart/form-data 格式
- 验证 dall-e-2/3 和 gpt-image-1 的 size 参数
- 设置默认 quality 和 size

#### `GetAndValidateEmbeddingRequest(c, relayMode) (*dto.EmbeddingRequest, error)`
- 验证 Embedding 请求

#### `GetAndValidateRerankRequest(c) (*dto.RerankRequest, error)`
- 验证 Rerank 请求，确保 query 和 documents 存在

#### `GetAndValidAudioRequest(c, relayMode) (*dto.AudioRequest, error)`
- 验证 Audio 请求

## 5. 关键逻辑分析

1. **模型默认值**: moderations 默认使用 `text-moderation-latest`，embeddings 从路径参数获取模型名
2. **图像大小验证**: dall-e-2 仅支持 256/512/1024，dall-e-3 支持 1024/1024x1792/1792x1024
3. **multipart 支持**: 图像编辑请求支持 multipart/form-data 格式
4. **FIM 支持**: chat completions 支持 prefix/suffix 字段（Fill-in-the-Middle）

## 6. 关联文件

- `dto/request.go`: 请求 DTO 定义
- `relay/constant/relay_mode.go`: RelayMode 常量
