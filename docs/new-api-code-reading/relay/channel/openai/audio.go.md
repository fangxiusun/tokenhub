# audio.go 代码阅读文档

## 1. 全局总结
该文件负责处理 OpenAI 的音频相关 API 响应，包括 TTS（文本转语音）和 STT（语音转文本）两种模式。TTS 处理器支持流式和非流式模式，非流式模式下会计算音频时长并据此估算 token 使用量。STT 处理器解析响应中的 usage 信息或基于输入长度估算。

## 2. 依赖关系
- `bytes`、`fmt`、`io`、`math`、`net/http` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/constant` — 常量定义
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `OpenaiTTSHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) *dto.Usage`
- **作用**：处理 OpenAI TTS（文本转语音）响应
- **参数**：Gin 上下文、HTTP 响应、中继信息
- **返回**：使用量信息
- **关键逻辑**：
  - 复制上游响应头到客户端
  - **流式模式**：使用 StreamScanner 逐块读取，检测包含 `usage` 的数据块提取 token 信息
  - **非流式模式**：
    1. 读取完整响应体并写入客户端
    2. 根据音频格式计算时长：
      - PCM 格式：直接根据字节数计算（24000Hz, 16-bit, 单声道）
      - 其他格式：使用 `common.GetAudioDuration` 获取时长
    3. 计算 CompletionTokens：`ceil(duration) / 60.0 * 1000`（每分钟 1000 tokens）
    4. 时长获取失败时的降级估算：`ceil(sizeInKB)` tokens

### `OpenaiSTTHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo, responseFormat string) (*types.NewAPIError, *dto.Usage)`
- **作用**：处理 OpenAI STT（语音转文本）响应
- **参数**：Gin 上下文、HTTP 响应、中继信息、响应格式
- **返回**：错误、使用量信息
- **关键逻辑**：
  1. 读取完整响应体并写入客户端
  2. 尝试从响应 JSON 中提取 usage 信息
  3. 如果响应包含有效的 usage，直接使用
  4. 否则基于输入长度估算 PromptTokens

## 5. 关键逻辑分析
- **音频时长计算**：PCM 格式无文件头，根据 OpenAI TTS 的固定参数（24000Hz/16-bit/单声道）直接从字节数计算时长
- **Token 估算策略**：当无法获取准确时长时，采用降级策略——每 KB 约等于 1 token
- **TTS Token 公式**：`ceil(duration_seconds) / 60.0 * 1000`，即每分钟 1000 tokens
- **Header 复制**：TTS 响应需要复制上游的 Content-Type 等头部（如 `audio/mpeg`）
- **流式 Usage 检测**：使用 `service.SundaySearch` 快速检测数据块中是否包含 "usage" 关键字

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到 TTS/STT 处理器
- `relay/helper/stream_scanner.go` — 流式扫描器
- `common/audio.go` — GetAudioDuration 音频时长获取
- `service/token.go` — GetEstimatePromptTokens 估算函数
