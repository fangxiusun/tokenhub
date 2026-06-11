# relay_adaptor.go 代码阅读文档

## 1. 全局总结

本文件是 Relay 模块的核心适配器工厂，负责根据 API 类型或任务平台返回对应的渠道适配器实例。文件包含两个工厂函数：`GetAdaptor` 用于同步请求（文本/图像/音频等），`GetTaskAdaptor` 用于异步任务（视频生成/音乐生成等）。

## 2. 依赖关系

- `constant`: API 类型常量、渠道类型常量、任务平台常量
- `relay/channel`: 所有渠道适配器包（ali, claude, gemini, openai 等 30+ 渠道）
- `relay/channel/task/*`: 任务适配器包（suno, kling, sora, vidu 等）
- `gin`: HTTP 框架

## 3. 类型定义

本文件无自定义类型定义，仅使用 `channel.Adaptor` 和 `channel.TaskAdaptor` 接口。

## 4. 函数详解

### `GetAdaptor(apiType int) channel.Adaptor`
- **功能**: 根据 API 类型返回对应的同步请求适配器
- **支持的渠道**: Ali, Anthropic, Baidu, Gemini, OpenAI, PaLM, Tencent, Xunfei, Zhipu, Ollama, Perplexity, AWS, Cohere, Dify, Jina, Cloudflare, SiliconFlow, Vertex, Mistral, DeepSeek, MokaAI, VolcEngine, BaiduV2, OpenRouter, Xinference, Xai, Coze, Jimeng, Moonshot, Submodel, MiniMax, Replicate, Codex
- **特殊处理**: OpenRouter 和 Xinference 复用 OpenAI 适配器
- **返回值**: 未匹配时返回 nil

### `GetTaskPlatform(c *gin.Context) constant.TaskPlatform`
- **功能**: 从请求上下文中提取任务平台标识
- **逻辑**: 优先读取 `channel_type` 整数字段，为空时读取 `platform` 字符串字段

### `GetTaskAdaptor(platform constant.TaskPlatform) channel.TaskAdaptor`
- **功能**: 根据任务平台返回对应的任务适配器
- **支持的平台**: Suno, Ali, Kling, Jimeng, Vertex, Vidu, Doubao/VolcEngine, Sora/OpenAI, Gemini, MiniMax/Hailuo
- **逻辑**: 先匹配字符串平台名（如 Suno），再尝试将平台标识解析为整数渠道类型进行匹配

## 5. 关键逻辑分析

1. **工厂模式**: 使用简单的 switch-case 工厂模式，每个 API 类型对应一个适配器实例
2. **复用策略**: OpenRouter 和 Xinference 共用 OpenAI 适配器，说明它们的 API 格式与 OpenAI 兼容
3. **任务适配器**: 任务适配器通过渠道类型数字进行匹配，支持多个渠道类型映射到同一个适配器（如 Doubao 和 VolcEngine 都映射到 taskdoubao）
4. **防御性编程**: 两个工厂函数在未匹配时返回 nil，调用方需进行空值检查

## 6. 关联文件

- `relay/channel/adapter.go`: 定义 `Adaptor` 和 `TaskAdaptor` 接口
- `relay/constant/relay_mode.go`: Relay 模式常量
- `constant/`: API 类型和渠道类型常量定义
