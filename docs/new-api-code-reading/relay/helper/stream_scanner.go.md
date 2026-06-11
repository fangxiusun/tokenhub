# stream_scanner.go 代码阅读文档

## 1. 全局总结

本文件实现了 SSE 流式响应的扫描和分发逻辑 `StreamScannerHandler`。是所有流式响应处理的核心引擎，负责从上游响应中读取 SSE 数据、分发给处理器、管理 Ping 保活、处理超时和错误。

## 2. 依赖关系

- `common`: 工具函数
- `relay/common`: StreamStatus
- `operation_setting`: 通用设置
- `bytedance/gopool`: 协程池

## 3. 类型定义

本文件无自定义类型定义，但定义了重要常量：
- `InitialScannerBufferSize`: 64KB 初始缓冲区
- `DefaultMaxScannerBufferSize`: 128MB 最大缓冲区
- `DefaultPingInterval`: 10 秒默认 Ping 间隔

## 4. 函数详解

### `NewStreamScanner(reader) *bufio.Scanner`
- 创建带缓冲区的流扫描器

### `StreamScannerHandler(c, resp, info, dataHandler)`
- **功能**: SSE 流处理主循环
- **架构**: 三个协程协作
  1. **Scanner 协程**: 读取上游响应，解析 SSE 数据，发送到 dataChan
  2. **Handler 协程**: 从 dataChan 读取数据，调用 dataHandler 处理
  3. **Ping 协程**: 定期发送 SSE ping 保活
- **主循环**: 等待超时、停止信号或客户端断开

## 5. 关键逻辑分析

1. **三协程架构**: Scanner → dataChan → Handler，Ping 独立运行
2. **超时控制**: 使用 ticker 重置超时，每个有效数据到达时重置
3. **Ping 保活**: 可配置的 ping 间隔，支持禁用
4. **资源清理**: defer 确保所有协程退出，最多等待 5 秒
5. **缓冲区管理**: 初始 64KB，最大 128MB（可配置）
6. **SSE 解析**: 过滤空行和非 data: 前缀的行，处理 [DONE] 终止信号

## 6. 关联文件

- `relay/helper/stream_result.go`: StreamResult 封装
- `relay/common/stream_status.go`: StreamStatus 状态管理
- `relay/helper/common.go`: SetEventStreamHeaders, PingData
