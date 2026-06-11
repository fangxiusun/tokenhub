# model.go 代码阅读文档

## 1. 全局总结

该文件实现了 AI 模型列表接口，兼容 OpenAI、Anthropic、Gemini 等多种 API 格式。在 `init()` 中构建全局模型列表，支持按用户组/Token 组过滤可用模型。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — API 类型、渠道类型常量
- `dto` — OpenAI/Anthropic/Gemini 模型结构
- `model` — 模型和渠道模型
- `relay` — 适配器获取
- `relay/channel/ai360/lingyiwanwu/minimax/moonshot` — 各渠道模型列表
- `relay/common` — RelayInfo
- `relay/helper` — 模型计费配置检查
- `service` — 自动组查询
- `setting/operation_setting` — 自用模式开关
- `types` — 错误类型
- `samber/lo` — 集合去重

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `modelListGroups` | 模型列表的组信息（userGroup、tokenGroup、ownerGroups） |

## 4. 函数详解

### `init()`
初始化全局模型列表。遍历所有 API 类型的适配器获取模型列表，加上 AI360、Moonshot、Lingyiwanwu、Minimax、Midjourney 等特殊渠道的模型。

### `ListModels(c *gin.Context, modelType int)`
核心模型列表接口。根据模型类型返回 OpenAI/Anthropic/Gemini 格式。支持 token 模型限制、组过滤、自动组。

### `ChannelListModels` / `DashboardListModels` / `EnabledListModels`
分别返回全部模型、渠道-模型映射、已启用模型。

### `RetrieveModel(c *gin.Context, modelType int)`
获取单个模型详情。

### `getPreferredModelOwners(modelNames, groups) map[string]string`
获取模型的首选所有者（渠道名称），用于显示模型来源。

### `buildOpenAIModel(modelName, ownerByModel) dto.OpenAIModels`
构建 OpenAI 格式的模型对象，注入端点支持信息。

## 5. 关键逻辑分析

- `init()` 在程序启动时构建全局模型列表和渠道-模型映射
- 模型过滤优先级：Token 模型限制 > 组启用模型
- `acceptUnsetRatioModel`（自用模式）允许显示未配置计费的模型
- 支持三种 API 格式输出：OpenAI、Anthropic、Gemini

## 6. 关联文件

- `relay/` — 各渠道适配器
- `model/channel.go` — 渠道模型和能力表
- `relay/helper/model.go` — 模型计费配置检查
