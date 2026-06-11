# context_key.go 代码阅读文档

## 1. 全局概述

本文件定义了 Gin Context 中使用的所有键名常量。这些键用于在请求生命周期的不同阶段（中间件 → 控制器 → 服务 → 模型）之间传递数据，是整个请求上下文的数据总线。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### ContextKey 类型

```go
type ContextKey string
```

自定义字符串类型，用于避免 Context 键名冲突。相比使用原生 `string` 类型，`ContextKey` 类型在编译时提供类型安全。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 上下文键分类

#### Token 计数相关
| 键名 | 用途 |
|------|------|
| `ContextKeyTokenCountMeta` | Token 计数元数据 |
| `ContextKeyPromptTokens` | Prompt Token 数量 |
| `ContextKeyEstimatedTokens` | 估算的 Token 数量 |
| `ContextKeyOriginalModel` | 原始模型名称（映射前） |
| `ContextKeyRequestStartTime` | 请求开始时间 |

#### Token 配额相关
| 键名 | 用途 |
|------|------|
| `ContextKeyTokenUnlimited` | 无限配额标记 |
| `ContextKeyTokenKey` | Token 密钥 |
| `ContextKeyTokenId` | Token ID |
| `ContextKeyTokenGroup` | Token 分组 |
| `ContextKeyTokenSpecificChannelId` | 指定渠道 ID |
| `ContextKeyTokenModelLimitEnabled` | 模型限制启用标记 |
| `ContextKeyTokenModelLimit` | 模型限制列表 |
| `ContextKeyTokenCrossGroupRetry` | 跨分组重试标记 |

#### 渠道相关
| 键名 | 用途 |
|------|------|
| `ContextKeyChannelId` | 渠道 ID |
| `ContextKeyChannelName` | 渠道名称 |
| `ContextKeyChannelCreateTime` | 渠道创建时间 |
| `ContextKeyChannelBaseUrl` | 渠道 Base URL |
| `ContextKeyChannelType` | 渠道类型 |
| `ContextKeyChannelSetting` | 渠道设置 |
| `ContextKeyChannelOtherSetting` | 渠道其他设置 |
| `ContextKeyChannelParamOverride` | 参数覆写 |
| `ContextKeyChannelHeaderOverride` | Header 覆写 |
| `ContextKeyChannelOrganization` | 组织 ID |
| `ContextKeyChannelAutoBan` | 自动封禁标记 |
| `ContextKeyChannelModelMapping` | 模型映射 |
| `ContextKeyChannelStatusCodeMapping` | 状态码映射 |
| `ContextKeyChannelIsMultiKey` | 多密钥标记 |
| `ContextKeyChannelMultiKeyIndex` | 多密钥索引 |
| `ContextKeyChannelKey` | 渠道密钥 |

#### 自动分组相关
| 键名 | 用途 |
|------|------|
| `ContextKeyAutoGroup` | 自动分组标记 |
| `ContextKeyAutoGroupIndex` | 自动分组索引 |
| `ContextKeyAutoGroupRetryIndex` | 自动分组重试索引 |

#### 用户相关
| 键名 | 用途 |
|------|------|
| `ContextKeyUserId` | 用户 ID（键名为 "id"） |
| `ContextKeyUserSetting` | 用户设置 |
| `ContextKeyUserQuota` | 用户配额 |
| `ContextKeyUserStatus` | 用户状态 |
| `ContextKeyUserEmail` | 用户邮箱 |
| `ContextKeyUserGroup` | 用户分组 |
| `ContextKeyUsingGroup` | 当前使用的分组 |
| `ContextKeyUserName` | 用户名 |

#### 其他
| 键名 | 用途 |
|------|------|
| `ContextKeyLocalCountTokens` | 本地 Token 计数标记 |
| `ContextKeySystemPromptOverride` | System Prompt 覆写 |
| `ContextKeyFileSourcesToCleanup` | 请求结束时需清理的文件源 |
| `ContextKeyAdminRejectReason` | 管理员拒绝原因（不返回给用户） |
| `ContextKeyLanguage` | 用户语言偏好（i18n） |
| `ContextKeyIsStream` | 是否流式请求 |

### 设计模式

所有键名使用 `ContextKey` 类型定义，而非原生 `string`，这是一种常见的 Go Context 安全模式，可以防止不同包之间的键名冲突。

## 6. 相关文件

- `middleware/` — 中间件设置这些 Context 值
- `controller/` — 控制器读取这些值
- `relay/` — 中继层使用渠道和 Token 相关的 Context 值
- `i18n/i18n.go` — 使用 `ContextKeyLanguage` 获取用户语言
