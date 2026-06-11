# channel-billing.go 代码阅读文档

## 1. 全局总结

该文件负责渠道（Channel）余额查询与更新。支持多种上游提供商（OpenAI、AIProxy、API2GPT、AIGC2D、SiliconFlow、DeepSeek、OpenRouter、Moonshot 等）的余额查询接口适配，以及批量更新所有渠道余额的功能。

## 2. 依赖关系

- `common` — 通用工具、渠道状态常量
- `constant` — 渠道类型常量
- `model` — 渠道数据模型
- `service` — 代理 HTTP 客户端、渠道禁用逻辑
- `setting/operation_setting` — 价格设置（Moonshot 汇率转换）
- `types` — 渠道错误类型
- `shopspring/decimal` — 精确十进制运算
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `OpenAISubscriptionResponse` | OpenAI 订阅信息响应 |
| `OpenAIUsageDailyCost` | OpenAI 每日费用 |
| `OpenAICreditGrants` | OpenAI 额度授予信息 |
| `OpenAIUsageResponse` | OpenAI 使用量响应（单位：0.01 美元） |
| `OpenAISBUsageResponse` | OpenAI-SB 使用量 |
| `AIProxyUserOverviewResponse` | AIProxy 用户概览 |
| `API2GPTUsageResponse` | API2GPT 使用量 |
| `APGC2DGPTUsageResponse` | AIGC2D 使用量 |
| `SiliconFlowUsageResponse` | SiliconFlow 用户信息 |
| `DeepSeekUsageResponse` | DeepSeek 余额信息 |
| `OpenRouterCreditResponse` | OpenRouter 额度 |

## 4. 函数详解

### `GetAuthHeader(token string) http.Header`
生成标准 Bearer Token 认证头。

### `GetClaudeAuthHeader(token string) http.Header`
生成 Claude（Anthropic）专用认证头（x-api-key + anthropic-version）。

### `GetResponseBody(method, url string, channel *model.Channel, headers http.Header) ([]byte, error)`
发送 HTTP 请求并返回响应体。支持渠道代理配置。

### `updateChannel*Balance(channel *model.Channel) (float64, error)`
各提供商余额查询函数：`updateChannelCloseAIBalance`、`updateChannelOpenAISBBalance`、`updateChannelAIProxyBalance`、`updateChannelAPI2GPTBalance`、`updateChannelSiliconFlowBalance`、`updateChannelDeepSeekBalance`、`updateChannelAIGC2DBalance`、`updateChannelOpenRouterBalance`、`updateChannelMoonshotBalance`。

### `updateChannelBalance(channel *model.Channel) (float64, error)`
根据渠道类型分发到对应的余额查询函数。OpenAI 类型还会额外查询 usage 接口计算实际余额。

### `UpdateChannelBalance(c *gin.Context)`
单个渠道余额更新的 HTTP 处理器。

### `updateAllChannelsBalance() error`
遍历所有已启用渠道，逐个更新余额。余额不足时自动禁用渠道。

### `UpdateAllChannelsBalance(c *gin.Context)`
批量更新所有渠道余额的 HTTP 处理器。

### `AutomaticallyUpdateChannels(frequency int)`
后台定时任务，周期性更新所有渠道余额。

## 5. 关键逻辑分析

- 多密钥渠道（`IsMultiKey`）不支持余额查询
- 余额不足（`balance <= 0`）时自动禁用渠道（`service.DisableChannel`）
- Moonshot 余额为 CNY，需除以 `operation_setting.Price` 转换为 USD
- 更新间隔通过 `common.RequestInterval` 控制，防止请求过快

## 6. 关联文件

- `controller/billing.go` — 用户级订阅/使用量接口
- `model/channel.go` — 渠道模型和余额更新方法
- `service/channel.go` — 渠道禁用/启用逻辑
