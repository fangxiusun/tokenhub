# operation_setting.go 代码阅读文档

## 1. 全局总结

该文件定义运营相关的全局开关和关键词列表，包括演示站点模式、自用模式、以及渠道自动禁用的关键词。

## 2. 依赖关系

- `strings` — 字符串处理

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `DemoSiteEnabled` | `bool` | 是否启用演示站点模式 |
| `SelfUseModeEnabled` | `bool` | 是否启用自用模式 |
| `AutomaticDisableKeywords` | `[]string` | 渠道自动禁用关键词列表 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `AutomaticDisableKeywordsToString` | `func AutomaticDisableKeywordsToString() string` | 将关键词列表转为换行分隔字符串 |
| `AutomaticDisableKeywordsFromString` | `func AutomaticDisableKeywordsFromString(s string)` | 从换行分隔字符串解析关键词列表 |

## 5. 关键逻辑分析

- 默认禁用关键词包括："Your credit balance is too low"、"This organization has been disabled."、"You exceeded your current quota" 等
- `AutomaticDisableKeywordsFromString` 会将所有关键词转为小写

## 6. 关联文件

- `relay/handler.go` — 使用自动禁用关键词
- `setting/operation_setting/status_code_ranges.go` — 状态码相关配置
