# suffix.go 代码阅读文档

## 1. 全局总结

该文件处理推理（Reasoning）模型的思考预算后缀解析，支持通用后缀、OpenAI 后缀和 DeepSeek V4 后缀。

## 2. 依赖关系

- `strings` — 字符串处理
- `github.com/samber/lo` — 集合操作工具

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `EffortSuffixes` | `[]string` | 通用思考预算后缀列表 |
| `OpenAIEffortSuffixes` | `[]string` | OpenAI 专用后缀列表 |
| `DeepSeekV4EffortSuffixes` | `[]string` | DeepSeek V4 专用后缀列表 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `TrimEffortSuffix` | `func TrimEffortSuffix(modelName string) (string, string, bool)` | 移除通用思考预算后缀 |
| `TrimEffortSuffixWithSuffixes` | `func TrimEffortSuffixWithSuffixes(modelName string, suffixes []string) (string, string, bool)` | 使用指定后缀列表移除后缀 |
| `ParseOpenAIReasoningEffortFromModelSuffix` | `func ParseOpenAIReasoningEffortFromModelSuffix(modelName string) (string, string)` | 解析 OpenAI 推理努力级别 |
| `ParseDeepSeekV4ThinkingSuffix` | `func ParseDeepSeekV4ThinkingSuffix(modelName string) (baseModel string, thinkingType string, effort string, ok bool)` | 解析 DeepSeek V4 思考后缀 |

## 5. 关键逻辑分析

- 通用后缀：`-max`、`-xhigh`、`-high`、`-medium`、`-low`、`-minimal`
- OpenAI 后缀额外包含 `-none`
- DeepSeek V4 仅支持 `-none`（禁用思考）和 `-max`（最大思考）
- `ParseDeepSeekV4ThinkingSuffix` 要求模型名以 `deepseek-v4-` 开头

## 6. 关联文件

- `relay/handler.go` — 使用后缀解析
- `relay/openai/` — OpenAI 适配器
- `relay/deepseek/` — DeepSeek 适配器
