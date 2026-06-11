# stream_status.go 代码阅读文档

## 1. 全局总结

本文件定义了流式响应状态管理结构 `StreamStatus`，用于跟踪 SSE 流的结束原因、错误记录和状态查询。

## 2. 依赖关系

- 标准库：sync, time, fmt, strings

## 3. 类型定义

### `StreamEndReason` (string)
流结束原因常量：
- `done`: 正常完成
- `timeout`: 超时
- `client_gone`: 客户端断开
- `scanner_error`: 扫描器错误
- `handler_stop`: 处理器主动停止
- `eof`: 到达流末尾
- `panic`: 协程 panic
- `ping_fail`: Ping 发送失败

### `StreamErrorEntry`
错误条目：Message, Timestamp

### `StreamStatus`
流状态结构：
- `EndReason`: 结束原因
- `EndError`: 结束错误
- `Errors`: 错误列表（最多 20 条）
- `ErrorCount`: 错误总数

## 4. 函数详解

### `NewStreamStatus() *StreamStatus`
- 创建新的流状态实例

### `SetEndReason(reason, err)`
- 设置结束原因（仅执行一次，sync.Once）

### `RecordError(msg)`
- 记录软错误（线程安全）

### `IsNormalEnd() bool`
- 判断是否正常结束（done/eof/handler_stop）

### `Summary() string`
- 生成状态摘要字符串

## 5. 关键逻辑分析

1. **线程安全**: 使用 sync.Once 和 sync.Mutex 保证并发安全
2. **错误限制**: 最多记录 20 条错误条目，避免内存无限增长
3. **正常结束判断**: done/eof/handler_stop 视为正常结束

## 6. 关联文件

- `relay/helper/stream_scanner.go`: 流扫描器使用 StreamStatus
- `relay/helper/stream_result.go`: StreamResult 封装
