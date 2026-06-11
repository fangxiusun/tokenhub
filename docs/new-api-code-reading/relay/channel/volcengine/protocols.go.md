# protocols.go 代码阅读文档

## 1. 全局总结
火山引擎 WebSocket 二进制协议的实现文件。定义了消息类型、事件类型、序列化/反序列化逻辑，用于与火山引擎 TTS 服务的 WebSocket 通信。

## 2. 依赖关系
- **标准库**: bytes, encoding/binary, fmt, io, math
- **外部依赖**: `github.com/gorilla/websocket`

## 3. 类型定义

### 协议位字段类型
```go
type (
    EventType         int32    // 事件类型
    MsgType           uint8    // 消息类型
    MsgTypeFlagBits   uint8    // 消息类型标志
    VersionBits       uint8    // 协议版本
    HeaderSizeBits    uint8    // 头部大小
    SerializationBits uint8    // 序列化方式
    CompressionBits   uint8    // 压缩方式
)
```

### 消息类型常量
```go
const (
    MsgTypeFullClientRequest    MsgType = 0b1      // 完整客户端请求
    MsgTypeAudioOnlyClient      MsgType = 0b10     // 仅音频客户端
    MsgTypeFullServerResponse   MsgType = 0b1001   // 完整服务端响应
    MsgTypeAudioOnlyServer      MsgType = 0b1011   // 仅音频服务端
    MsgTypeFrontEndResultServer MsgType = 0b1100   // 前端结果
    MsgTypeError                MsgType = 0b1111   // 错误消息
)
```

### 事件类型常量
覆盖连接管理、会话管理、TTS、ASR、Chat、字幕等多种事件类型。

### Message 结构体
```go
type Message struct {
    Version       VersionBits
    HeaderSize    HeaderSizeBits
    MsgType       MsgType
    MsgTypeFlag   MsgTypeFlagBits
    Serialization SerializationBits
    Compression   CompressionBits
    EventType     EventType
    SessionID     string
    ConnectID     string
    Sequence      int32
    ErrorCode     uint32
    Payload       []byte
}
```
二进制协议消息结构体，支持序列化和反序列化。

## 4. 函数详解

### NewMessageFromBytes(data []byte) (*Message, error)
从字节数组解析消息，先解析消息类型和标志，再调用 `Unmarshal`。

### NewMessage(msgType MsgType, flag MsgTypeFlagBits) (*Message, error)
创建新消息实例，设置默认值（Version1、HeaderSize4、JSON 序列化、无压缩）。

### (m *Message) Marshal() ([]byte, error)
将消息序列化为字节数组：
1. 构建 3 字节头部（版本+头部大小、类型+标志、序列化+压缩）
2. 按消息类型写入可选字段（事件、会话ID、序列号、错误码）
3. 写入载荷

### (m *Message) Unmarshal(data []byte) error
从字节数组反序列化消息，按头部大小跳过填充字节，然后按类型读取可选字段。

### ReceiveMessage(conn *websocket.Conn) (*Message, error)
从 WebSocket 连接接收消息并反序列化。

### FullClientRequest(conn *websocket.Conn, payload []byte) error
发送完整的客户端请求消息。

## 5. 关键逻辑分析

1. **二进制协议**: 使用大端序（BigEndian）的二进制协议，头部固定 3 字节（版本、类型、序列化/压缩），可选头部大小为 4 字节的倍数。

2. **消息类型分发**: 不同消息类型有不同的可选字段：
   - 普通消息: 可能包含序列号
   - 事件消息: 包含事件类型和会话ID
   - 错误消息: 包含错误码

3. **会话连接管理**: 连接相关的事件（StartConnection、ConnectionStarted 等）不包含 SessionID，而会话相关的事件需要 SessionID。

4. **TTS 流式通信**: 通过 `MsgTypeAudioOnlyServer` 消息传输音频数据，负序列号表示最后一包。

## 6. 关联文件
- `volcengine/tts.go` — 使用 `FullClientRequest` 和 `ReceiveMessage` 进行 WebSocket 通信
- `volcengine/adaptor.go` — 在 `handleTTSWebSocketResponse` 中使用
