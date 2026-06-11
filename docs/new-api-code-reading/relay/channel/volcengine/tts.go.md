# tts.go 代码阅读文档

## 1. 全局总结
火山引擎 TTS（文本转语音）功能的实现文件。包含 TTS 请求/响应数据结构、语音类型映射、编码格式映射、HTTP 和 WebSocket 两种响应处理方式。

## 2. 依赖关系
- **标准库**: context, encoding/base64, encoding/json, errors, fmt, io, net/http, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/dto` — Usage 结构体
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/google/uuid` — UUID 生成
  - `github.com/gorilla/websocket` — WebSocket 客户端

## 3. 类型定义

### TTS 请求结构体
```go
type VolcengineTTSRequest struct {
    App     VolcengineTTSApp     `json:"app"`
    User    VolcengineTTSUser    `json:"user"`
    Audio   VolcengineTTSAudio   `json:"audio"`
    Request VolcengineTTSReqInfo `json:"request"`
}
```
- **App**: 应用信息（AppID、Token、Cluster）
- **User**: 用户信息（UID）
- **Audio**: 音频配置（语音类型、编码、语速、采样率等）
- **Request**: 请求信息（请求ID、文本、操作类型、模型等）

### TTS 响应结构体
```go
type VolcengineTTSResponse struct {
    ReqID    string                     `json:"reqid"`
    Code     int                        `json:"code"`
    Message  string                     `json:"message"`
    Sequence int                        `json:"sequence"`
    Data     string                     `json:"data"`
    Addition *VolcengineTTSAdditionInfo `json:"addition,omitempty"`
}
```

### 映射表
- `openAIToVolcengineVoiceMap`: OpenAI 语音名 → 火山引擎语音类型
- `responseFormatToEncodingMap`: 响应格式 → 编码格式

## 4. 函数详解

### parseVolcengineAuth(apiKey string) (appID, token string, err error)
解析 API Key，格式为 `appid|access_token`。

### mapVoiceType(openAIVoice string) string
将 OpenAI 语音名映射为火山引擎语音类型。

### mapEncoding(responseFormat string) string
将响应格式映射为编码格式。

### getContentTypeByEncoding(encoding string) string
根据编码格式返回对应的 Content-Type。

### handleTTSResponse(...)
非流式 TTS 响应处理：
1. 读取响应体并解析为 `VolcengineTTSResponse`
2. 检查响应码（3000 表示成功）
3. Base64 解码音频数据
4. 设置 Content-Type 并返回音频

### handleTTSWebSocketResponse(...)
流式 TTS WebSocket 响应处理：
1. 建立 WebSocket 连接
2. 发送 `FullClientRequest`
3. 循环接收消息：
   - `MsgTypeAudioOnlyServer`: 写入音频数据
   - 负序列号: 表示最后一包
   - `MsgTypeError`: 返回错误
4. 设置 Transfer-Encoding 为 chunked

### generateRequestID() string
生成 UUID 作为请求 ID。

## 5. 关键逻辑分析

1. **双协议支持**: 火山引擎 TTS 同时支持 HTTP（非流式）和 WebSocket（流式）两种协议。

2. **音频编码映射**: OpenAI 的 `opus` 格式映射为火山引擎的 `ogg_opus`，`aac` 和 `flac` 映射为 `mp3`。

3. **WebSocket 通信流程**: 
   - 建立连接 → 发送请求 → 接收音频包 → 检测结束标志
   - 结束标志: 序列号为负数

4. **Base64 音频数据**: 非流式响应的音频数据以 Base64 编码返回，需要解码后传递给客户端。

5. **情绪支持**: Audio 配置支持情绪参数（enable_emotion、emotion、emotion_scale），但需要模型支持。

## 6. 关联文件
- `volcengine/protocols.go` — WebSocket 二进制协议实现
- `volcengine/adaptor.go` — 调用 TTS 处理函数
- `volcengine/constants.go` — 模型列表（包含 TTS 模型）
