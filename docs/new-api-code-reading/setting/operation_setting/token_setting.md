# token_setting.go 代码阅读文档

## 1. 全局总结

该文件定义令牌（Token）相关配置，控制每用户最大令牌数量。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `TokenSetting` | `MaxUserTokens` | `int` | `1000` | 每用户最大令牌数量 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetTokenSetting` | `func GetTokenSetting() *TokenSetting` | 获取令牌配置 |
| `GetMaxUserTokens` | `func GetMaxUserTokens() int` | 获取每用户最大令牌数量 |

## 5. 关键逻辑分析

- 默认每用户最多 1000 个令牌

## 6. 关联文件

- `model/token.go` — 令牌模型
- `controller/token.go` — 令牌接口
