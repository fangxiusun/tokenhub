# model_sync.go 代码阅读文档

## 1. 全局总结

该文件实现了从上游（basellm.github.io）同步模型和供应商元数据的功能。支持创建缺失模型、选择性覆盖更新、差异预览，以及 ETag 缓存和重试机制。

## 2. 依赖关系

- `common` — 通用工具、HTTP 配置
- `model` — 模型和供应商模型
- `gin-gonic/gin` — HTTP 框架
- `gorm.io/gorm` — 数据库事务

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `upstreamEnvelope[T]` | 上游响应信封（泛型） |
| `upstreamModel` | 上游模型数据 |
| `upstreamVendor` | 上游供应商数据 |
| `overwriteField` | 覆盖字段指定 |
| `syncRequest` | 同步请求体 |

## 4. 函数详解

### `SyncUpstreamModels(c *gin.Context)`
同步上游模型和供应商。流程：获取缺失模型 → 拉取上游数据 → 创建缺失模型 → 处理可选覆盖更新。

### `SyncUpstreamPreview(c *gin.Context)`
预览上游与本地差异。返回缺失模型列表和冲突字段详情。

### `fetchJSON[T](ctx, url, out) error`
带重试和 ETag 缓存的 JSON 获取函数。支持指数退避和抖动。

### `ensureVendorID(vendorName, ...) 确保供应商存在，返回 ID。

## 5. 关键逻辑分析

- 上游数据源：`basellm.github.io/llm-metadata`
- 支持国际化（en、zh-CN、zh-TW、ja）
- ETag 缓存避免重复下载（304 Not Modified）
- 重试机制：默认 3 次，指数退避 + 随机抖动
- 覆盖更新需指定字段列表（description、icon、tags、vendor、name_rule、status）
- `sync_official = 0` 的模型不参与同步
- HTTP 客户端使用自定义 Dialer 优先尝试 IPv4 连接 github.io

## 6. 关联文件

- `model/model.go` — 模型元数据模型
- `model/vendor.go` — 供应商模型
- `controller/missing_models.go` — 缺失模型查询
