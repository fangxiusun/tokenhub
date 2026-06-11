# stream_scanner_test.go 代码阅读文档

## 1. 全局总结

本文件是 `stream_scanner.go` 的单元测试，全面测试了 SSE 流扫描器的正确性、解耦性、Ping 保活和 StreamStatus 集成。

## 2. 依赖关系

- `gin`: 测试上下文
- `relay/common`: StreamStatus
- `operation_setting`: 通用设置

## 3. 测试用例分类

### 基础正确性测试
- `TestStreamScannerHandler_NilInputs`: nil 输入安全
- `TestStreamScannerHandler_EmptyBody`: 空响应体
- `TestStreamScannerHandler_1000Chunks/10000Chunks`: 大量 chunk 处理
- `TestStreamScannerHandler_OrderPreserved`: 顺序保持
- `TestStreamScannerHandler_DoneStopsScanner`: [DONE] 停止扫描
- `TestStreamScannerHandler_StopStopsStream`: Stop 停止流
- `TestStreamScannerHandler_SkipsNonDataLines`: 跳过非 data 行
- `TestStreamScannerHandler_DataWithExtraSpaces`: 空格处理

### 解耦性测试
- `TestStreamScannerHandler_ScannerDecoupledFromSlowHandler`: Scanner 与 Handler 解耦
- `TestStreamScannerHandler_SlowUpstreamFastHandler`: 慢上游 + 快处理器

### Ping 测试
- `TestStreamScannerHandler_PingSentDuringSlowUpstream`: 慢上游时发送 Ping
- `TestStreamScannerHandler_PingDisabledByRelayInfo`: DisablePing 禁用 Ping
- `TestStreamScannerHandler_PingInterleavesWithSlowUpstream`: Ping 与数据交错

### StreamStatus 集成测试
- `TestStreamScannerHandler_StreamStatus_DoneReason`: Done 原因
- `TestStreamScannerHandler_StreamStatus_EOFWithoutDone`: EOF 原因
- `TestStreamScannerHandler_StreamStatus_HandlerStop`: Handler Stop 原因
- `TestStreamScannerHandler_StreamStatus_HandlerDone`: Handler Done 原因
- `TestStreamScannerHandler_StreamStatus_Timeout`: 超时原因
- `TestStreamScannerHandler_StreamStatus_SoftErrors`: 软错误
- `TestStreamScannerHandler_StreamStatus_MultipleErrorsPerChunk`: 每 chunk 多个错误
- `TestStreamScannerHandler_StreamStatus_ErrorThenStop`: 先错误后停止

## 4. 关键逻辑分析

1. **解耦验证**: 通过时间比较验证 Scanner 和 Handler 的解耦程度
2. **Ping 验证**: 通过响应体中 `: PING` 的出现次数验证 Ping 保活
3. **超时测试**: 通过修改全局 StreamingTimeout 测试超时行为
4. **并发安全**: 使用 atomic 操作和 sync.Mutex 测试并发场景

## 5. 关联文件

- `relay/helper/stream_scanner.go`: 被测试的源文件
- `relay/common/stream_status.go`: StreamStatus 类型
