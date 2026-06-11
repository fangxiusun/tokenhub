# websocket.go 代码阅读文档

## 1. 全局总结

本文件实现了 WebSocket（实时语音/对话）请求的处理入口 `WssHelper`。处理 WebSocket 连接的建立、消息转发和计费。

## 2. 依赖关系

- `relay/common`: RelayInfo
- `service`: 计费
- `gorilla/websocket`: WebSocket 客户端库

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `WssHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 WebSocket 实时请求
- **流程**: InitChannelMeta → 获取适配器 → DoRequest（建立 WS 连接） → DoResponse（消息转发循环） → PostWssConsumeQuota
- **特殊点**: 请求体为 nil，因为 WebSocket 连接通过 header 建立

## 5. 关键逻辑分析

1. **连接管理**: 响应体是 `*websocket.Conn`，存储在 `info.TargetWs` 中
2. **延迟关闭**: 使用 `defer info.TargetWs.Close()` 确保连接关闭
3. **实时计费**: 使用 `PostWssConsumeQuota` 进行实时使用量计费
4. **状态码映射**: 支持 statusCodeMapping 重置错误状态码

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor.DoRequest/DoResponse 接口
- `dto/realtime.go`: RealtimeUsage DTO
