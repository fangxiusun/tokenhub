# channel-test.go 代码阅读文档

## 1. 全局总结

该文件实现了渠道测试功能，包括单渠道测试、全渠道测试和自动定时测试。通过模拟真实请求来验证渠道的可用性、响应时间和模型支持情况，支持多种端点类型（Chat、Embedding、Image、Rerank、Responses 等）和流式/非流式模式。

## 2. 依赖关系

- `common` — 通用工具、日志
- `constant` — 渠道类型、端点类型常量
- `dto` — 请求/响应数据结构
- `middleware` — 上下文设置
- `model` — 渠道/日志模型
- `pkg/billingexpr` — 计费表达式
- `relay` — 适配器获取
- `relay/common` — RelayInfo 生成
- `relay/helper` — 模型映射/价格辅助
- `service` — 错误处理、渠道管理
- `setting/ratio_setting` — 模型后缀配置
- `types` — 价格数据、错误类型
- `bytedance/gopkg/util/gopool` — 协程池
- `samber/lo` — 集合操作
- `tidwall/gjson` — JSON 解析

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `testResult` | 测试结果，包含 gin.Context、本地错误和 NewAPI 错误 |

## 4. 函数详解

### `testChannel(channel, testUserID, testModel, endpointType, isStream) testResult`
核心测试函数。流程：构建测试请求 → 设置上下文 → 模型映射 → 价格计算 → 转换请求 → 发送请求 → 处理响应 → 结算配额 → 记录日志。

### `buildTestRequest(model, endpointType, channel, isStream) dto.Request`
根据端点类型和模型名构建测试请求。支持 Embedding、Image、Rerank、Responses、ResponsesCompact 等类型。

### `TestChannel(c *gin.Context)`
单渠道测试的 HTTP 处理器。支持从缓存获取渠道信息。

### `testAllChannels(notify bool) error`
全渠道测试。使用 gopool 协程池异步执行，支持禁用阈值检测和自动禁用/启用。

### `TestAllChannels(c *gin.Context)`
全渠道测试的 HTTP 处理器。

### `AutomaticallyTestChannels()`
自动定时测试的后台任务。仅在 Master 节点运行，根据 `AutoTestChannelMinutes` 配置的间隔周期执行。

### `settleTestQuota(info, priceData, usage) (int, *billingexpr.TieredResult)`
测试配额结算。支持分层计费和传统计费两种模式。

### `coerceTestUsage(usageAny, isStream, estimatePromptTokens) (*dto.Usage, error)`
类型强制转换 usage 对象，流式模式下使用估算的 prompt tokens。

## 5. 关键逻辑分析

- 不支持测试的渠道类型：Midjourney、SunoAPI、Kling、Jimeng、DoubaoVideo、Vidu
- 自动检测模型类型：rerank、embedding、codex、responses compaction
- 全渠道测试使用互斥锁防止重复运行
- 测试失败超过阈值或响应超时会自动禁用渠道
- 测试用户 ID 优先从上下文获取，否则查找 root 用户

## 6. 关联文件

- `controller/channel.go` — 渠道管理
- `relay/` — 适配器和请求转换
- `service/channel.go` — 渠道禁用/启用逻辑
- `middleware/context.go` — `SetupContextForSelectedChannel`
