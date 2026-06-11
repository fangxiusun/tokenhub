# OAuth 层详细设计 (`oauth/`)

## 1. 概述
OAuth 层为多个 OAuth 提供商（GitHub, Discord, OIDC 等）提供统一接口，使用注册表模式进行动态提供商管理。

## 2. 文件详细说明

### 2.1 `provider.go` - 核心接口
**职责**: 定义所有 OAuth 提供商必须实现的核心 `Provider` 接口。

**接口方法**:
- **`GetName()`**: 返回提供商的显示名称（如 "GitHub", "Discord"）。
- **`IsEnabled()`**: 返回此 OAuth 提供商是否启用。
- **`ExchangeToken(ctx context.Context, code string, c *gin.Context)`**: 将授权码交换为访问令牌。gin.Context 用于需要请求信息的提供商（如 redirect_uri）。
- **`GetUserInfo(ctx context.Context, token *OAuthToken)`**: 使用访问令牌检索用户信息。
- **`IsUserIDTaken(providerUserID string)`**: 检查提供商用户 ID 是否已与账户关联。
- **`FillUserByProviderID(user *model.User, providerUserID string)`**: 通过提供商用户 ID 填充用户模型。
- **`SetProviderUserID(user *model.User, providerUserID string)`**: 在用户模型上设置提供商用户 ID。
- **`GetProviderPrefix()`**: 返回自动生成用户名的前缀（如 "github_"）。

### 2.2 `registry.go` - 提供商注册表
**职责**: 提供商注册表，用于动态提供商管理。

**关键函数**:
- **`Register(name string, provider Provider)`**: 注册提供商。
- **`Get(name string)`**: 按名称检索提供商。
- **`GetAll()`**: 返回所有已注册的提供商。
- **`GetEnabled()`**: 返回所有已启用的提供商。

### 2.3 `types.go` - 共享类型
**职责**: OAuth 提供商的共享类型。

**关键结构体**:
- **`OAuthToken`**: 表示 OAuth 访问令牌。
- **`OAuthUser`**: 表示来自 OAuth 提供商的用户信息。

### 2.4 提供商实现
- **`github.go`** -- GitHub OAuth 提供商实现。
- **`discord.go`** -- Discord OAuth 提供商实现。
- **`linuxdo.go`** -- LinuxDO OAuth 提供商实现。
- **`oidc.go`** -- OpenID Connect (OIDC) 提供商实现。
- **`generic.go`** -- 自定义提供商的通用 OAuth 提供商实现。

## 3. 架构模式
OAuth 层使用 **注册表模式**：
1. 所有提供商都实现 `Provider` 接口。
2. 提供商在初始化期间向全局注册表注册自身。
3. 控制器层从注册表中按名称检索提供商。
4. 这允许动态添加/删除 OAuth 提供商，无需更改核心代码。

## 4. 使用流程
1. **初始化**: 在 `main.go` 中调用 `oauth.Init()` 注册所有提供商。
2. **授权**: 控制器调用 `oauth.Get(providerName)` 获取提供商，然后调用 `ExchangeToken()` 交换令牌。
3. **用户信息**: 使用访问令牌调用 `GetUserInfo()` 获取用户信息。
4. **绑定**: 调用 `FillUserByProviderID()` 和 `SetProviderUserID()` 处理用户绑定。

---

## 关联文件列表

### OAuth 层核心文件
- `oauth/provider.go` - Provider 接口定义
- `oauth/registry.go` - 提供商注册表
- `oauth/types.go` - 共享类型
- `oauth/github.go` - GitHub 提供商
- `oauth/discord.go` - Discord 提供商
- `oauth/linuxdo.go` - LinuxDO 提供商
- `oauth/oidc.go` - OIDC 提供商
- `oauth/generic.go` - 通用提供商

### 依赖的模型层文件
- `model/user.go` - 用户模型
- `model/user_oauth_binding.go` - 用户 OAuth 绑定模型

### 使用此模块的文件
- `main.go` - 应用入口（初始化 OAuth）
- `controller/oauth.go` - OAuth 统一处理
- `controller/github.go` - GitHub OAuth 控制器
- `controller/discord.go` - Discord OAuth 控制器
- `controller/linuxdo.go` - LinuxDO OAuth 控制器
- `controller/oidc.go` - OIDC 控制器
- `controller/custom_oauth.go` - 自定义 OAuth 控制器
