# custom_oauth.go 代码阅读文档

## 1. 全局总结

该文件实现了自定义 OAuth 提供商的完整 CRUD 管理，以及用户 OAuth 绑定的查看和解绑功能。支持通过 OIDC Discovery 自动发现端点配置。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 自定义 OAuth 提供商和绑定模型
- `oauth` — OAuth 提供商注册/注销
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `CustomOAuthProviderResponse` | OAuth 提供商响应（不含 client_secret） |
| `UserOAuthBindingResponse` | 用户 OAuth 绑定响应 |
| `CreateCustomOAuthProviderRequest` | 创建请求 |
| `UpdateCustomOAuthProviderRequest` | 更新请求 |
| `FetchCustomOAuthDiscoveryRequest` | Discovery 获取请求 |

## 4. 函数详解

### CRUD 接口
- `GetCustomOAuthProviders` — 获取所有提供商
- `GetCustomOAuthProvider` — 获取单个提供商
- `CreateCustomOAuthProvider` — 创建提供商（检查 slug 冲突）
- `UpdateCustomOAuthProvider` — 更新提供商（支持部分更新）
- `DeleteCustomOAuthProvider` — 删除提供商（检查用户绑定）

### Discovery
- `FetchCustomOAuthDiscovery` — 通过后端代理获取 OIDC Discovery 文档（避免 CORS）

### 用户绑定
- `GetUserOAuthBindings` — 获取当前用户的 OAuth 绑定
- `GetUserOAuthBindingsByAdmin` — 管理员查看指定用户的绑定
- `UnbindCustomOAuth` — 解绑当前用户的 OAuth
- `UnbindCustomOAuthByAdmin` — 管理员解绑指定用户的 OAuth

## 5. 关键逻辑分析

- 响应中排除 `client_secret` 字段，确保安全
- 创建/更新时检查 slug 是否与内置 OAuth 提供商冲突
- 删除时检查是否有用户绑定，有绑定则禁止删除
- slug 变更时自动注销旧 slug 并注册新 slug
- 支持字段级别的部分更新（`*string`、`*bool`、`*int` 指针类型）

## 6. 关联文件

- `oauth/` — OAuth 提供商注册机制
- `model/custom_oauth.go` — 自定义 OAuth 数据模型
- `controller/misc.go` — `GetStatus` 中注入自定义 OAuth 信息
