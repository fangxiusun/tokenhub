# realtime.go 代码阅读文档

## 1. 全局摘要

该文件定义了 OpenAI Realtime API 的事件和会话数据结构。Realtime API 用于实时语音交互，支持音频流、转录、工具调用等功能。

## 2. 依赖

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`OpenAIError` 类型

## 3. 类型定义

### 事件类型常量

**客户端发送事件**：
- `RealtimeEventTypeError`："error"
- `RealtimeEventTypeSessionUpdate`："session.update"
- `RealtimeEventTypeConversationCreate`："conversation.item.create"
- `RealtimeEventTypeResponseCreate`："response.create"
- `RealtimeEventInputAudioBufferAppend`："input_audio_buffer.append"

**服务器发送事件**：
- `RealtimeEventTypeResponseDone`："response.done"
- `RealtimeEventTypeSessionUpdated`："session.updated"
- `RealtimeEventTypeSessionCreated`："session.created"
- `RealtimeEventResponseAudioDelta`："response.audio.delta"
- `RealtimeEventResponseAudioTranscriptionDelta`："response.audio_transcript.delta"
- `RealtimeEventResponseFunctionCallArgumentsDelta`："response.function_call_arguments.delta"
- `RealtimeEventResponseFunctionCallArgumentsDone`："response.function_call_arguments.done"
- `RealtimeEventConversationItemCreated`："conversation.item.created"

### 核心结构体

**RealtimeEvent**：实时事件结构：
- `EventId` (string)：事件 ID
- `Type` (string)：事件类型
- `Session` (*RealtimeSession)：会话信息
- `Item` (*RealtimeItem)：会话项
- `Error` (*types.OpenAIError)：错误信息
- `Response` (*RealtimeResponse)：响应信息
- `Delta` (string)：增量数据
- `Audio` (string)：音频数据（Base64）

**RealtimeResponse**：实时响应：
- `Usage` (*RealtimeUsage)：使用量统计

**RealtimeUsage**：使用量统计：
- `TotalTokens` (int)：总 token 数
- `InputTokens` (int)：输入 token 数
- `OutputTokens` (int)：输出 token 数
- `InputTokenDetails` (InputTokenDetails)：输入 token 详情
- `OutputTokenDetails` (OutputTokenDetails)：输出 token 详情

### 会话相关结构体

**RealtimeSession**：实时会话：
- `Modalities` ([]string)：模态列表
- `Instructions` (string)：指令
- `Voice` (string)：语音类型
- `InputAudioFormat` (string)：输入音频格式
- `OutputAudioFormat` (string)：输出音频格式
- `InputAudioTranscription` (InputAudioTranscription)：输入音频转录配置
- `TurnDetection` (interface{})：轮次检测配置
- `Tools` ([]RealTimeTool)：工具列表
- `ToolChoice` (string)：工具选择
- `Temperature` (float64)：温度参数

**InputAudioTranscription**：音频转录配置：
- `Model` (string)：转录模型

**RealTimeTool**：实时工具：
- `Type` (string)：工具类型
- `Name` (string)：工具名称
- `Description` (string)：工具描述
- `Parameters` (any)：参数定义

### 会话项结构体

**RealtimeItem**：实时会话项：
- `Id` (string)：项 ID
- `Type` (string)：项类型
- `Status` (string)：状态
- `Role` (string)：角色
- `Content` ([]RealtimeContent)：内容数组
- `Name` (*string)：名称
- `ToolCalls` (any)：工具调用
- `CallId` (string)：调用 ID

**RealtimeContent**：实时内容：
- `Type` (string)：内容类型
- `Text` (string)：文本内容
- `Audio` (string)：音频数据（Base64）
- `Transcript` (string)：转录文本

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **事件驱动架构**：Realtime API 使用事件驱动模型，客户端和服务器通过事件进行通信。

2. **音频流处理**：支持实时音频流传输，音频数据使用 Base64 编码。

3. **语音转录**：支持实时语音转录，将音频转换为文本。

4. **工具调用**：支持实时工具调用，用于扩展功能。

5. **会话管理**：通过 `RealtimeSession` 管理会话状态和配置。

## 6. 相关文件

- `relay/realtime/`：Realtime 中继适配器
- `controller/realtime.go`：Realtime 控制器
- `types/realtime.go`：Realtime 类型定义