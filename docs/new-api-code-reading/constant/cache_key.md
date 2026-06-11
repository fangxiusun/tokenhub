# cache_key.go 代码阅读文档

## 1. 全局概述

本文件定义了 Redis 缓存键的格式常量，用于统一管理系统中的缓存键命名规范。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 用户相关缓存键格式

| 常量名 | 格式 | 用途 |
|--------|------|------|
| `UserGroupKeyFmt` | `user_group:%d` | 缓存用户所属分组，`%d` 为用户 ID |
| `UserQuotaKeyFmt` | `user_quota:%d` | 缓存用户配额，`%d` 为用户 ID |
| `UserEnabledKeyFmt` | `user_enabled:%d` | 缓存用户启用状态，`%d` 为用户 ID |
| `UserUsernameKeyFmt` | `user_name:%d` | 缓存用户名，`%d` 为用户 ID |

### Token 相关字段名

| 常量名 | 值 | 用途 |
|--------|-----|------|
| `TokenFiledRemainQuota` | `RemainQuota` | Token 剩余额度字段名 |
| `TokenFieldGroup` | `Group` | Token 分组字段名 |

### 设计说明

- 使用 `%d` 格式化占位符，便于通过 `fmt.Sprintf` 生成具体的缓存键
- 缓存键命名采用 `前缀:标识` 的规范格式
- Token 字段名用于数据库查询和缓存操作中的字段引用

## 6. 相关文件

- `common/redis.go` — Redis 操作封装
- `model/user.go` — 用户模型，使用这些缓存键
- `model/token.go` — Token 模型，使用 Token 字段常量
