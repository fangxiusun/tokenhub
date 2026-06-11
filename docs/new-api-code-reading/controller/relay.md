# relay.go 代码阅读文档

## 1. 全局总结

该文件是 AI API 中继（Relay）的核心入口。处理所有 AI 请求的转发，包括文本、图像、音频、嵌入、Rerank、Responses 等类型。支持重试、渠道自动禁用、计费预扣和错误日志记录。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — 常量
- `dto` — 数据结构
- `logger` — 日志
- `middleware` — 上下文设置
- `model` — 渠道/日志模型
- `pkg/perf_metrics` — 性能指标
- `relay` — 各类型 Helper
- `relay/common` — RelayInfo
- `relay/constant` — RelayMode 常量
- `relay/helper` — 价格/模型映射辅助
- `service` — 渠道选择、计费、错误处理
- `setting` — 敏感词检查
- `setting/operation_setting` — 重试配置
- `types` — 错误类型
- `bytedance/gopkg/util/gopool` — 协程池
- `gorilla/websocket` — WebSocket

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `Relay(c *gin.Context, relayFormat types.RelayFormat)`
核心中继函数。流程：解析请求 → 生成 RelayInfo → 敏感词检查 → Token 估算 → 价格计算 → 预扣费 → 重试循环 → 调用对应 Handler → 错误处理和退费。

### `relayHandler(c, info)` / `geminiRelayHandler(c, info)`
根据 RelayMode 分发到对应的 Helper 函数。

### `RelayMidjourney(c *gin.Context)`
Midjourney 专用中继。支持通知、任务查询、图像种子、换脸等操作。

### `RelayTask(c *gin.Context)` / `RelayTaskFetch(c *gin.Context)`
Task 类型中继（视频生成等长任务）。支持锁定渠道和重试。

### `getChannel(c, info, retryParam)` / `shouldRetry(c, err, retryTimes)`
渠道选择和重试判断。

### `processChannelError(c, channelError, err)`
处理渠道错误：自动禁用、记录错误日志。

## 5. 关键逻辑分析

- 支持 WebSocket（Realtime API）连接升级
- 重试时会重新获取渠道（可能切换到不同渠道）
- 预扣费在请求失败时自动退费
- 错误日志包含完整的渠道信息和重试历史
- Task 类型支持锁定渠道（避免重试时切换渠道）
- 敏感词检查和 Token 估算可根据配置跳过

## 6. 关联文件

- `relay/` — 各类型 Helper 和适配器
- `service/channel.go` — 渠道选择
- `service/billing.go` — 计费
- `middleware/` — 上下文设置
