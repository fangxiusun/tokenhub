# relay_realtime.go 代码阅读文档

## 1. 全局总结
该文件实现了 OpenAI Realtime API（实时语音/文本对话）的 WebSocket 中继处理。通过双向 WebSocket 代理，将客户端连接与上游 OpenAI Realtime 服务桥接，同时进行实时 token 计数和配额预扣费。支持文本和音频两种 token 类型的分别计数。

## 2. 依赖关系
- `fmt` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — WebSocket 辅助函数
- `github.com/QuantumNous/new-api/service` — 业务逻辑（token 计数、配额预扣）
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/bytedance/gopkg/util/gopool` — 协程池
- `github.com/gin-gonic/gin` — HTTP 框架
- `github.com/gorilla/websocket` — WebSocket 库

## 3. 类型定义
无自定义类型定义。使用 `dto.RealtimeEvent`、`dto.RealtimeUsage` 等外部 DTO。

## 4. 函数详解

### `OpenaiRealtimeHandler(c *gin.Context, info *relaycommon.RelayInfo) (*types.NewAPIError, *dto.RealtimeUsage)`
- **作用**：处理 Realtime API 的 WebSocket 连接
- **参数**：Gin 上下文、中继信息（包含 ClientWs 和 TargetWs）
- **返回**：错误、累计使用量
- **关键逻辑**：
  1. 启动两个 goroutine：客户端读取器和目标读取器
  2. **客户端读取器**：
     - 读取客户端 WebSocket 消息
     - 解析 RealtimeEvent，提取 session update 中的 tools 配置
     - 计算输入 token（text + audio）
     - 转发消息到上游
  3. **目标读取器**：
     - 读取上游 WebSocket 消息
     - 根据事件类型分发处理：
       - `response.done` — 累加 usage 并预扣费，重置 usage
       - `session.updated/created` — 更新音频格式信息
       - 其他事件 — 计算输出 token 并累加
     - 转发消息到客户端
  4. 等待任一连接关闭或错误
  5. 处理剩余未计费的 usage

### `preConsumeUsage(ctx *gin.Context, info *relaycommon.RelayInfo, usage *dto.RealtimeUsage, totalUsage *dto.RealtimeUsage) error`
- **作用**：累加使用量并执行预扣费
- **逻辑**：将本次 usage 累加到 totalUsage，然后调用 `service.PreWssConsumeQuota` 扣费

## 5. 关键逻辑分析
- **双向代理**：使用 goroutine + channel 模式实现全双工 WebSocket 代理，sendChan/receiveChan 用于可选的消息录制
- **实时计费**：每个 `response.done` 事件触发一次预扣费，确保长连接的配额实时消耗
- **Token 分类**：分别统计 text tokens 和 audio tokens，支持详细的使用量报告
- **工具配置跟踪**：从 `session.update` 事件中提取 tools 配置，存储在 `info.RealtimeTools` 中
- **音频格式更新**：从 `session.updated/created` 事件中更新输入/输出音频格式
- **Fallback 计费**：当上游未返回 usage 时，通过 `service.CountTokenRealtime` 本地估算
- **连接关闭处理**：通过 channel 关闭信号（clientClosed/targetClosed）优雅处理连接断开
- **Panic 恢复**：每个 goroutine 都有 defer recover，防止 panic 导致整个处理崩溃

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到 Realtime 处理器
- `relay/helper/websocket.go` — WssString WebSocket 发送辅助
- `service/realtime.go` — CountTokenRealtime、PreWssConsumeQuota
- `dto/realtime.go` — RealtimeEvent、RealtimeUsage DTO 定义
- `relay/common/relay_info.go` — ClientWs、TargetWs 字段
