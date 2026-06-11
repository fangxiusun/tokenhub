# keys.go 代码阅读文档

## 1. 全局概述

本文件定义了后端国际化（i18n）系统的所有消息键常量。这些常量用于替代硬编码的字符串，确保所有用户可见的消息都支持多语言翻译。消息键按业务模块分类组织。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

本文件无自定义类型定义，所有常量为字符串类型。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 消息键分类

#### 通用错误消息（common）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgInvalidParams` | `common.invalid_params` | 无效参数 |
| `MsgDatabaseError` | `common.database_error` | 数据库错误 |
| `MsgRetryLater` | `common.retry_later` | 请稍后重试 |
| `MsgGenerateFailed` | `common.generate_failed` | 生成失败 |
| `MsgNotFound` | `common.not_found` | 未找到 |
| `MsgUnauthorized` | `common.unauthorized` | 未授权 |
| `MsgForbidden` | `common.forbidden` | 禁止访问 |
| `MsgOperationSuccess` | `common.operation_success` | 操作成功 |
| `MsgOperationFailed` | `common.operation_failed` | 操作失败 |
| `MsgAlreadyExists` | `common.already_exists` | 已存在 |
| `MsgBatchTooMany` | `common.batch_too_many` | 批量操作过多 |

#### 认证中间件消息（auth）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgAuthNotLoggedIn` | `auth.not_logged_in` | 未登录 |
| `MsgAuthAccessTokenInvalid` | `auth.access_token_invalid` | 访问令牌无效 |
| `MsgAuthUserBanned` | `auth.user_banned` | 用户已被封禁 |
| `MsgAuthInsufficientPrivilege` | `auth.insufficient_privilege` | 权限不足 |

#### Token 相关消息（token）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgTokenExpired` | `token.expired` | Token 已过期 |
| `MsgTokenExhausted` | `token.exhausted` | Token 额度已用尽 |
| `MsgTokenInvalid` | `token.invalid` | Token 无效 |

#### 兑换码相关消息（redemption）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgRedemptionUsed` | `redemption.used` | 兑换码已使用 |
| `MsgRedemptionExpired` | `redemption.expired` | 兑换码已过期 |

#### 用户相关消息（user）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgUserExists` | `user.exists` | 用户已存在 |
| `MsgUserDisabled` | `user.disabled` | 用户已被禁用 |
| `MsgUserTransferSuccess` | `user.transfer_success` | 转账成功 |

#### 配额相关消息（quota）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgQuotaInsufficient` | `quota.insufficient` | 配额不足 |

#### 订阅相关消息（subscription）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgSubscriptionNotEnabled` | `subscription.not_enabled` | 订阅未启用 |

#### 支付相关消息（payment）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgPaymentNotConfigured` | `payment.not_configured` | 支付未配置 |

#### 渠道相关消息（channel）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgChannelNotExists` | `channel.not_exists` | 渠道不存在 |
| `MsgChannelNoAvailableKey` | `channel.no_available_key` | 无可用密钥 |

#### 模型相关消息（model）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgModelNameExists` | `model.name_exists` | 模型名称已存在 |

#### OAuth 相关消息（oauth）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgOAuthInvalidCode` | `oauth.invalid_code` | OAuth 授权码无效 |
| `MsgOAuthGetUserErr` | `oauth.get_user_error` | 获取用户信息失败 |
| `MsgOAuthTrustLevelLow` | `oauth.trust_level_low` | 信任等级过低 |

#### 分发器相关消息（distributor）
| 常量名 | 键值 | 说明 |
|--------|------|------|
| `MsgDistributorNoAvailableChannel` | `distributor.no_available_channel` | 无可用渠道 |

### 命名规范

消息键采用 `模块.具体错误` 的格式：
- 模块名使用小写英文
- 具体错误使用下划线分隔的英文短语
- 常量名使用 `Msg` 前缀 + 模块名 + 错误描述（驼峰命名）

### 设计原则

1. **统一管理**：所有用户可见的消息都通过消息键管理
2. **类型安全**：使用常量而非字符串字面量，编译时检查
3. **模块化组织**：按业务模块分组，便于维护
4. **扩展性**：新增消息只需添加常量和对应的翻译文件

## 6. 相关文件

- `i18n/i18n.go` — i18n 初始化和翻译函数
- `i18n/locales/zh-CN.yaml` — 中文翻译文件
- `i18n/locales/en.yaml` — 英文翻译文件
- `i18n/locales/zh-TW.yaml` — 繁体中文翻译文件
- `controller/` — 控制器中使用这些消息键
