# model_ratio.go 代码阅读文档

## 1. 全局总结

该文件是整个倍率系统的核心，定义了所有模型的输入倍率、补全倍率、价格、图像倍率、音频倍率等配置，支持通配符匹配和思考预算模型处理。

## 2. 依赖关系

- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/common` — 系统日志、JSON 操作
- `github.com/QuantumNous/new-api/setting/operation_setting` — 运营设置
- `github.com/QuantumNous/new-api/types` — RWMap 并发安全 Map

## 3. 类型定义

| 常量名 | 值 | 说明 |
|--------|------|------|
| `USD2RMB` | `7.3` | 美元兑人民币汇率 |
| `USD` | `500` | $0.002 = 1 额度 |
| `RMB` | `500/7.3` | 人民币换算系数 |

| 结构体 | 字段 | 说明 |
|--------|------|------|
| `CompletionRatioInfo` | `Ratio float64` | 补全倍率 |
| | `Locked bool` | 是否锁定（不可修改） |

主要映射表：
- `defaultModelRatio` — 模型输入倍率（283+ 条）
- `defaultModelPrice` — 模型固定价格
- `defaultCompletionRatio` — 补全倍率
- `defaultCacheRatio` — 缓存倍率
- `defaultCreateCacheRatio` — 缓存创建倍率
- `defaultImageRatio` — 图像倍率
- `defaultAudioRatio` — 音频倍率
- `defaultAudioCompletionRatio` — 音频补全倍率

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `InitRatioSettings` | `func InitRatioSettings()` | 初始化所有倍率设置 |
| `GetModelPrice` | `func GetModelPrice(name string, printErr bool) (float64, bool)` | 获取模型价格 |
| `GetModelRatio` | `func GetModelRatio(name string) (float64, bool, string)` | 获取模型输入倍率 |
| `GetCompletionRatio` | `func GetCompletionRatio(name string) float64` | 获取模型补全倍率 |
| `GetCompletionRatioInfo` | `func GetCompletionRatioInfo(name string) CompletionRatioInfo` | 获取补全倍率信息（含锁定状态） |
| `GetModelRatioOrPrice` | `func GetModelRatioOrPrice(model string) (float64, bool, bool)` | 获取模型倍率或价格 |
| `FormatMatchingModelName` | `func FormatMatchingModelName(name string) string` | 格式化模型名用于匹配 |
| `GetAudioRatio` | `func GetAudioRatio(name string) float64` | 获取音频倍率 |
| `GetAudioCompletionRatio` | `func GetAudioCompletionRatio(name string) float64` | 获取音频补全倍率 |
| `GetImageRatio` | `func GetImageRatio(name string) (float64, bool)` | 获取图像倍率 |

## 5. 关键逻辑分析

- 1 额度 = $0.002 = ¥0.014
- `FormatMatchingModelName` 处理思考预算模型：将 `gemini-2.5-flash-preview-xxx-thinking-yyy` 映射为 `gemini-2.5-flash-thinking-*`
- `getHardcodedCompletionModelRatio` 使用硬编码逻辑处理特殊模型的补全倍率
- GPT-5 系列补全倍率为 8，GPT-4o 系列为 4
- Claude 系列补全倍率为 5
- 未找到的模型默认返回 37.5 倍率

## 6. 关联文件

- `setting/ratio_setting/cache_ratio.go` — 缓存倍率
- `setting/ratio_setting/group_ratio.go` — 组倍率
- `service/billing.go` — 计费服务
- `relay/handler.go` — 使用倍率计算费用
