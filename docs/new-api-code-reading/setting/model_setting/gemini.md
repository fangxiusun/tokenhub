# gemini.go 代码阅读文档

## 1. 全局总结

该文件定义 Gemini 模型的配置，包括安全设置、版本设置、图像生成支持模型列表、Thinking 适配器等。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `GeminiSettings` | `SafetySettings` | `map[string]string` | 安全设置（默认 "OFF"） |
| | `VersionSettings` | `map[string]string` | API 版本设置（默认 "v1beta"） |
| | `SupportedImagineModels` | `[]string` | 支持图像生成的模型列表 |
| | `ThinkingAdapterEnabled` | `bool` | 是否启用 Thinking 适配器 |
| | `ThinkingAdapterBudgetTokensPercentage` | `float64` | Thinking 预算百分比 |
| | `FunctionCallThoughtSignatureEnabled` | `bool` | 是否启用函数调用思考签名 |
| | `RemoveFunctionResponseIdEnabled` | `bool` | 是否移除函数响应 ID |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetGeminiSettings` | `func GetGeminiSettings() *GeminiSettings` | 获取 Gemini 配置 |
| `GetGeminiSafetySetting` | `func GetGeminiSafetySetting(key string) string` | 获取安全设置 |
| `GetGeminiVersionSetting` | `func GetGeminiVersionSetting(key string) string` | 获取版本设置 |
| `IsGeminiModelSupportImagine` | `func IsGeminiModelSupportImagine(model string) bool` | 检查模型是否支持图像生成 |

## 5. 关键逻辑分析

- 安全设置默认 "OFF"，允许自定义每个模型的安全级别
- 版本设置默认 "v1beta"，`gemini-1.0-pro` 使用 "v1"
- 注意：bool 字段必须以 `enabled` 结尾才能在编辑界面生效

## 6. 关联文件

- `relay/gemini/` — Gemini 中继适配器
- `setting/model_setting/global.go` — 全局模型设置
