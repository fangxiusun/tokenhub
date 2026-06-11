# ratio_sync.go 代码阅读文档

## 1. 全局总结

该文件实现了从多个上游源（自定义上游、OpenRouter、models.dev）同步倍率配置的功能。支持差异对比、并发获取、格式转换和可信度评估。

## 2. 依赖关系

- `common` — 通用工具
- `dto` — 上游请求/差异数据结构
- `logger` — 日志
- `model` — 渠道模型
- `setting/billing_setting` — 计费同步字段
- `setting/ratio_setting` — 倍率配置
- `gin-gonic/gin` — HTTP 框架
- `samber/lo` — 集合操作

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `upstreamResult` | 上游获取结果 |
| `modelsDevProvider` | models.dev 供应商结构 |
| `modelsDevModel` | models.dev 模型结构 |
| `modelsDevCost` | models.dev 成本结构 |
| `modelsDevCandidate` | models.dev 候选定价 |

## 4. 函数详解

### `FetchUpstreamRatios(c *gin.Context)`
核心同步函数。支持四种上游格式：
- Type1: `/api/ratio_config` — 静态倍率配置
- Type2: `/api/pricing` — 定价列表
- Type3: OpenRouter `/v1/models` — 按 token 计费
- Type4: models.dev `/api.json` — 供应商级定价

### `GetSyncableChannels(c *gin.Context)`
获取可同步的渠道列表，包含官方预设和 models.dev 预设。

### `convertOpenRouterToRatioData(reader)`
将 OpenRouter 的每 token USD 定价转换为本地倍率格式。

### `convertModelsDevToRatioData(reader)`
将 models.dev 的每百万 token USD 定价转换为本地倍率格式。

## 5. 关键逻辑分析

- 并发获取：最多 8 个并发，支持重试（3 次，指数退避）
- 可信度评估：model_ratio=37.5 且 completion_ratio=1 的数据标记为不可信
- 官方预设 ID：-100，models.dev 预设 ID：-101
- models.dev 冲突解决：选择最低非零 input cost
- HTTP 客户端优先尝试 IPv4 连接 github.io

## 6. 关联文件

- `dto/upstream.go` — 上游请求/结果结构
- `setting/ratio_setting/` — 倍率配置
- `setting/billing_setting/` — 计费同步
