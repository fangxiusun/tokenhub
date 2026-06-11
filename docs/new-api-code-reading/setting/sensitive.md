# sensitive.go 代码阅读文档

## 1. 全局总结

该文件管理敏感词检测配置，支持 Prompt 和 Completion 的敏感词检查，以及流式模式下的缓存队列配置。

## 2. 依赖关系

- `strings` — 字符串处理

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `CheckSensitiveEnabled` | `bool` | `true` | 是否启用敏感词检测 |
| `CheckSensitiveOnPromptEnabled` | `bool` | `true` | 是否在 Prompt 中检测敏感词 |
| `StopOnSensitiveEnabled` | `bool` | `true` | 检测到敏感词时是否停止生成 |
| `StreamCacheQueueLength` | `int` | `0` | 流模式缓存队列长度 |
| `SensitiveWords` | `[]string` | `["test_sensitive"]` | 敏感词列表 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `SensitiveWordsToString` | `func SensitiveWordsToString() string` | 将敏感词列表转为换行分隔的字符串 |
| `SensitiveWordsFromString` | `func SensitiveWordsFromString(s string)` | 从换行分隔的字符串解析敏感词列表 |
| `ShouldCheckPromptSensitive` | `func ShouldCheckPromptSensitive() bool` | 判断是否应检查 Prompt 敏感词 |

## 5. 关键逻辑分析

- `SensitiveWordsFromString` 会先清空列表，然后按换行符分割并去除空行和首尾空白
- `ShouldCheckPromptSensitive` 需要同时满足 `CheckSensitiveEnabled` 和 `CheckSensitiveOnPromptEnabled`

## 6. 关联文件

- `middleware/sensitive.go` — 敏感词检测中间件
- `service/sensitive.go` — 敏感词服务
