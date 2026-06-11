# tools.go 代码阅读文档

## 1. 全局总结

该文件管理工具调用（Tool Call）定价配置，支持按模型前缀设置不同价格，并包含 GPT Image 1 和 Gemini 音频输入的特殊定价逻辑。

## 2. 依赖关系

- `sort` — 排序
- `strings` — 字符串处理
- `sync/atomic` — 原子操作
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `ToolPriceSetting` | `Prices` | `map[string]float64` | 工具价格映射 |
| `prefixEntry` | `prefix` | `string` | 模型前缀 |
| | `price` | `float64` | 价格 |
| `toolPriceIndex` | `defaults` | `map[string]float64` | 默认价格 |
| | `prefixes` | `map[string][]prefixEntry` | 按工具名的前缀价格列表 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `RebuildToolPriceIndex` | `func RebuildToolPriceIndex()` | 重建价格查找索引 |
| `GetToolPriceForModel` | `func GetToolPriceForModel(toolName, modelName string) float64` | 获取工具的模型特定价格 |
| `GetToolPrice` | `func GetToolPrice(toolName string) float64` | 获取工具的默认价格 |
| `GetGPTImage1PriceOnceCall` | `func GetGPTImage1PriceOnceCall(quality string, size string) float64` | 获取 GPT Image 1 单次调用价格 |
| `GetGeminiInputAudioPricePerMillionTokens` | `func GetGeminiInputAudioPricePerMillionTokens(modelName string) float64` | 获取 Gemini 音频输入价格 |

## 5. 关键逻辑分析

- 价格查找顺序：最长前缀匹配 → 工具默认价格 → 硬编码回退 → 0
- 键格式：`"tool_name"` 或 `"tool_name:model_prefix*"`
- 使用 `atomic.Pointer` 实现无锁读取
- GPT Image 1 定价按质量（low/medium/high）和尺寸分档
- Gemini 音频定价按模型版本分档

## 6. 关联文件

- `relay/handler.go` — 使用工具价格
- `service/billing.go` — 计费服务
