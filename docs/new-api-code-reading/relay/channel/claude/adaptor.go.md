# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是 Claude 渠道的适配器实现，负责将 OpenAI 格式的请求转换为 Claude API 格式。支持 Claude 原生 API 格式，处理请求头、URL 构建和响应格式转换。

## 2. 依赖关系
- 标准库：`errors`, `fmt`, `io`, `net/http`, `net/url`
- 内部包：
  - `dto`: 数据传输对象
  - `relay/channel`: 渠道通用工具
  - `relaycommon`: 中继通用配置
  - `model_setting`: 模型设置
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
### 结构体
- `Adaptor`: Claude 渠道适配器结构体（空结构体）。

## 4. 函数详解
### 核心函数
1. **`ConvertClaudeRequest`**: 直接返回原始 Claude 请求，不做转换。
2. **`GetRequestURL`**: 构建 Claude API URL，支持 Beta 查询参数。
3. **`shouldAppendClaudeBetaQuery`**: 判断是否需要添加 Beta 查询参数。
4. **`CommonClaudeHeadersOperation`**: 处理 Claude 通用请求头，包括 anthropic-beta 和模型设置头。
5. **`SetupRequestHeader`**: 设置请求头，包括 API Key、Anthropic 版本和通用头。
6. **`ConvertOpenAIRequest`**: 将 OpenAI 请求转换为 Claude 格式。
7. **`DoRequest`**: 执行 API 请求。
8. **`DoResponse`**: 处理 Claude 响应，区分流式和非流式模式。

## 5. 关键逻辑分析
- **原生格式支持**：Claude 请求直接传递，不做格式转换。
- **Beta 功能支持**：通过查询参数 `beta=true` 启用 Beta 功能。
- **请求头管理**：设置 `x-api-key`、`anthropic-version` 等必需头。
- **模型设置头**：通过 `model_setting.GetClaudeSettings().WriteHeaders` 应用模型特定的头设置。
- **流式/非流式处理**：根据 `info.IsStream` 标志选择不同的响应处理器。

## 6. 关联文件
- `claude/constants.go`: 定义模型列表和渠道名称。
- `claude/dto.go`: Claude 特定的数据传输对象（已注释）。
- `claude/relay-claude.go`: 请求转换和响应处理逻辑。
- `relay/channel/openai/adaptor.go`: OpenAI 适配器，用于 OpenAI 请求转换。
- `model_setting/claude.go`: Claude 模型设置。