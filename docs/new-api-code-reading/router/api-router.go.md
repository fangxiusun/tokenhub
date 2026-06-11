# api-router.go 代码阅读文档

## 1. 全局总结

`router/api-router.go` 是系统最核心的路由文件，定义了 `/api` 路径下的所有 REST API 端点。涵盖用户认证、渠道管理、Token 管理、日志查询、订阅计费、OAuth 集成、部署管理等完整业务功能。约 400 行代码，注册了 150+ 个路由端点。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/QuantumNous/new-api/controller` | 所有请求处理器 |
| `github.com/QuantumNous/new-api/middleware` | 认证、限流、缓存等中间件 |
| `github.com/QuantumNous/new-api/oauth` | OAuth 提供商注册（通过 init()） |
| `github.com/gin-contrib/gzip` | Gzip 压缩中间件 |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

本文件无自定义类型，主要使用 Gin 的 `RouterGroup` 进行路由分组。

## 4. 函数详解

### SetApiRouter(router *gin.Engine)

**功能**：注册所有 `/api` 前缀下的 REST API 路由。

**全局中间件**（应用于所有 `/api` 路由）：
- `middleware.RouteTag("api")` — 路由标签标记
- `gzip.Gzip(gzip.DefaultCompression)` — 响应压缩
- `middleware.BodyStorageCleanup()` — 请求体存储清理
- `middleware.GlobalAPIRateLimit()` — 全局 API 限流

## 5. 关键逻辑分析

### 路由分组结构

#### 公开端点（无需认证）
| 路由 | 功能 |
|------|------|
| `GET /setup` | 检查系统初始化状态 |
| `POST /setup` | 系统初始化（带匿名限流） |
| `GET /status` | 系统状态 |
| `GET /notice` | 公告信息 |
| `GET /user-agreement` | 用户协议 |
| `GET /privacy-policy` | 隐私政策 |
| `GET /home_page_content` | 首页内容 |
| `GET /pricing` | 价格配置（需 pricing 模块权限） |
| `GET /rankings` | 排行榜（需 rankings 模块权限） |
| `GET /ratio_config` | 比率配置 |

#### 用户认证路由 `/api/user`
| 路由 | 功能 |
|------|------|
| `POST /register` | 注册（带 Turnstile 验证） |
| `POST /login` | 登录（带 Turnstile 验证） |
| `POST /login/2fa` | 两步验证码登录 |
| `POST /passkey/login/begin` | Passkey 登录开始 |
| `POST /passkey/login/finish` | Passkey 登录完成 |
| `GET /logout` | 登出 |
| `POST /epay/notify` | 支付回调 |

#### 用户自管理路由 `/api/user/self`（需 UserAuth）
| 路由 | 功能 |
|------|------|
| `GET /self` | 获取当前用户信息 |
| `PUT /self` | 更新个人信息 |
| `DELETE /self` | 删除账号 |
| `GET /token` | 生成 API Token |
| `GET /passkey` | Passkey 状态 |
| `POST /passkey/register/*` | Passkey 注册流程 |
| `POST /passkey/verify/*` | Passkey 验证流程 |
| `DELETE /passkey` | 删除 Passkey |
| `GET /aff` | 获取邀请码 |
| `POST /topup` | 充值 |
| `POST /pay` | 支付请求 |
| `POST /stripe/pay` | Stripe 支付 |
| `POST /creem/pay` | Creem 支付 |
| `POST /waffo/pay` | Waffo 支付 |
| `PUT /setting` | 更新用户设置 |
| `GET/POST /2fa/*` | 两步验证管理 |
| `GET/POST /checkin` | 签到功能 |
| `GET/DELETE /oauth/bindings` | OAuth 绑定管理 |

#### 管理员用户管理 `/api/user`（需 AdminAuth）
| 路由 | 功能 |
|------|------|
| `GET /` | 获取所有用户 |
| `GET /search` | 搜索用户 |
| `POST /` | 创建用户 |
| `POST /manage` | 管理用户 |
| `PUT /` | 更新用户 |
| `DELETE /:id` | 删除用户 |
| `GET /topup` | 所有充值记录 |
| `POST /topup/complete` | 完成充值 |
| `DELETE /:id/reset_passkey` | 重置用户 Passkey |
| `GET /2fa/stats` | 2FA 统计 |
| `DELETE /:id/2fa` | 禁用用户 2FA |

#### 订阅计费 `/api/subscription`（需 UserAuth）
| 路由 | 功能 |
|------|------|
| `GET /plans` | 获取订阅计划 |
| `GET /self` | 获取当前订阅 |
| `PUT /self/preference` | 更新订阅偏好 |
| `POST /balance/pay` | 余额支付 |
| `POST /epay/pay` | Epay 支付 |
| `POST /stripe/pay` | Stripe 支付 |
| `POST /creem/pay` | Creem 支付 |

#### 订阅管理 `/api/subscription/admin`（需 AdminAuth）
| 路由 | 功能 |
|------|------|
| `GET/POST /plans` | 订阅计划 CRUD |
| `PUT /plans/:id` | 更新计划 |
| `PATCH /plans/:id` | 更新计划状态 |
| `POST /bind` | 绑定订阅 |
| `GET/POST /users/:id/subscriptions` | 用户订阅管理 |

#### 系统配置 `/api/option`（需 RootAuth）
| 路由 | 功能 |
|------|------|
| `GET /` | 获取系统选项 |
| `PUT /` | 更新系统选项 |
| `POST /payment_compliance` | 确认支付合规 |
| `POST /rest_model_ratio` | 重置模型比率 |
| `POST /waffo-pancake/*` | Waffo Pancake 集成 |

#### 自定义 OAuth 提供商 `/api/custom-oauth-provider`（需 RootAuth）
| 路由 | 功能 |
|------|------|
| `POST /discovery` | OAuth 发现 |
| `GET/POST /` | 提供商 CRUD |
| `PUT/DELETE /:id` | 更新/删除提供商 |

#### 性能管理 `/api/performance`（需 RootAuth）
| 路由 | 功能 |
|------|------|
| `GET /stats` | 性能统计 |
| `DELETE /disk_cache` | 清理磁盘缓存 |
| `POST /reset_stats` | 重置统计 |
| `POST /gc` | 强制 GC |
| `GET/DELETE /logs` | 日志管理 |

#### 渠道管理 `/api/channel`（需 AdminAuth）
| 路由 | 功能 |
|------|------|
| `GET/POST /` | 渠道 CRUD |
| `GET /search` | 搜索渠道 |
| `GET /models` | 渠道模型列表 |
| `POST /:id/key` | 获取渠道密钥（需 RootAuth + 安全验证） |
| `GET /test` | 测试所有渠道 |
| `GET /test/:id` | 测试单个渠道 |
| `GET /update_balance` | 更新渠道余额 |
| `POST /tag/*` | 标签管理 |
| `POST /batch` | 批量删除 |
| `POST /fix` | 修复渠道能力 |
| `POST /fetch_models` | 获取上游模型 |
| `POST /codex/*` | Codex OAuth 管理 |
| `POST /ollama/*` | Ollama 模型管理 |
| `POST /copy/:id` | 复制渠道 |
| `POST /multi_key/manage` | 多密钥管理 |
| `POST /upstream_updates/*` | 上游模型更新检测与应用 |

#### Token 管理 `/api/token`（需 UserAuth）
| 路由 | 功能 |
|------|------|
| `GET /` | 获取所有 Token |
| `GET /search` | 搜索 Token |
| `POST /` | 创建 Token |
| `PUT /` | 更新 Token |
| `DELETE /:id` | 删除 Token |
| `POST /batch` | 批量删除 |
| `POST /batch/keys` | 批量获取密钥 |

#### 日志管理 `/api/log`
| 路由 | 权限 | 功能 |
|------|------|------|
| `GET /` | Admin | 所有日志 |
| `DELETE /` | Admin | 清理历史日志 |
| `GET /stat` | Admin | 日志统计 |
| `GET /self/stat` | User | 自己的日志统计 |
| `GET /search` | Admin | 搜索日志 |
| `GET /self` | User | 自己的日志 |
| `GET /self/search` | User | 搜索自己的日志 |

#### 兑换码管理 `/api/redemption`（需 AdminAuth）
标准 CRUD 操作，包含搜索和批量删除无效兑换码。

#### 其他管理路由
| 路由组 | 功能 |
|--------|------|
| `/api/group` | 用户组管理 |
| `/api/prefill_group` | 预填充组管理 |
| `/api/mj` | Midjourney 任务查询 |
| `/api/task` | 任务查询 |
| `/api/vendors` | 供应商元数据管理 |
| `/api/models` | 模型元数据管理（含上游同步） |
| `/api/deployments` | 模型部署管理（IoNet 集成） |

### 中间件使用模式

- **认证中间件**：`UserAuth()` → `AdminAuth()` → `RootAuth()`（权限递增）
- **限流中间件**：`CriticalRateLimit()` 用于敏感操作，`GlobalAPIRateLimit()` 用于全局
- **安全中间件**：`TurnstileCheck()` 人机验证，`SecureVerificationRequired()` 安全验证
- **缓存中间件**：`DisableCache()` 用于需要实时数据的端点
- **模块权限**：`HeaderNavModuleAuth("pricing")` 按模块控制访问

### Webhook 端点

系统支持多个支付平台的 Webhook 回调：
- `POST /api/stripe/webhook` — Stripe 支付回调
- `POST /api/creem/webhook` — Creem 支付回调
- `POST /api/waffo/webhook` — Waffo 支付回调
- `POST /api/waffo-pancake/webhook/:env` — Waffo Pancake 回调（按环境分离）

### OAuth 集成

- **标准 OAuth**：`GET /api/oauth/:provider`（GitHub、Discord、OIDC、LinuxDO）
- **非标准 OAuth**：微信（`/api/oauth/wechat`）、Telegram（`/api/oauth/telegram/login`）
- **邮箱绑定**：`POST /api/oauth/email/bind`

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/main.go` | 调用 SetApiRouter 注册路由 |
| `controller/` | 所有路由的处理器实现 |
| `middleware/` | 认证、限流、缓存中间件 |
| `oauth/` | OAuth 提供商注册（通过 import side-effect） |
| `constant/` | API 类型、渠道类型常量 |
| `types/` | 请求/响应类型定义 |
