# constant.go 代码阅读文档

## 1. 全局总结
该文件定义了 OpenAI 渠道的常量，包含一个非常庞大的模型列表，涵盖了 OpenAI 发布的所有模型系列（GPT-3.5、GPT-4、GPT-4o、GPT-5、o1/o3/o4 推理模型、嵌入模型、图像生成模型、音频模型、实时模型等）以及渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` (var []string)
- **作用**：定义 OpenAI 渠道支持的完整模型列表
- **涵盖系列**：
  - **GPT-3.5 系列**：gpt-3.5-turbo 及其变体
  - **GPT-4 系列**：gpt-4、gpt-4-turbo、gpt-4-vision-preview
  - **GPT-4o 系列**：gpt-4o、gpt-4o-mini 及其变体（transcribe、search、tts）
  - **GPT-4.1 系列**：gpt-4.1、gpt-4.1-mini、gpt-4.1-nano
  - **GPT-4.5 系列**：gpt-4.5-preview
  - **GPT-5 系列**：gpt-5、gpt-5-mini、gpt-5-nano、gpt-5-pro 等
  - **GPT-5.1/5.2/5.3/5.4 系列**：最新迭代版本
  - **推理模型**：o1、o1-mini、o1-pro、o3、o3-mini、o3-pro、o4-mini
  - **深度研究**：o3-deep-research、o4-mini-deep-research
  - **嵌入模型**：text-embedding-ada-002、text-embedding-3-small/large
  - **图像生成**：dall-e-2、dall-e-3、gpt-image-1/1-mini/1.5、sora-2
  - **音频模型**：gpt-audio、gpt-audio-mini、gpt-realtime 系列
  - **语音模型**：whisper-1、tts-1 系列
  - **审核模型**：text-moderation、omni-moderation
  - **其他**：computer-use-preview、davinci-002、babbage-002

### `ChannelName` (var string)
- **作用**：渠道名称标识
- **值**：`"openai"`

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 模型列表非常全面，包含日期后缀的精确版本（如 `gpt-4o-2024-08-06`）和通用版本
- 涵盖了 OpenAI 从传统聊天模型到推理模型、多模态模型的完整产品线
- 部分模型有多个日期版本，支持精确版本选择

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — 通过 `GetModelList()` 引用此常量
- `relay/channel/openai/helper.go` — 流式响应处理中使用模型信息
