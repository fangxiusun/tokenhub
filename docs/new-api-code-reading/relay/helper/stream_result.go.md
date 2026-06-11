# stream_result.go 代码阅读文档

## 1. 全局总结

本文件定义了 `StreamResult` 结构，封装了流式处理过程中的状态控制逻辑，包括记录软错误、信号停止和标记完成。

## 2. 依赖关系

- `relay/common`: StreamStatus

## 3. 类型定义

### `StreamResult`
```go
type StreamResult struct {
    status  *relaycommon.StreamStatus
    stopped bool
}
```

## 4. 函数详解

### `Error(err)`
- 记录软错误，流继续处理

### `Stop(err)`
- 记录致命错误并标记流停止

### `Done()`
- 标记正常完成（如 Dify 的 message_end）

### `IsStopped() bool`
- 返回当前 chunk 是否已停止

### `reset()`
- 清除 per-chunk 停止标志

## 5. 关键逻辑分析

1. **双层控制**: Error 记录错误但不停止，Stop/Done 记录错误并停止
2. **Chunk 级别**: stopped 标志在每个 chunk 处理后重置
3. **状态委托**: 实际状态存储在 StreamStatus 中，StreamResult 只是便捷封装

## 6. 关联文件

- `relay/common/stream_status.go`: StreamStatus 实现
- `relay/helper/stream_scanner.go`: StreamScannerHandler 使用 StreamResult
