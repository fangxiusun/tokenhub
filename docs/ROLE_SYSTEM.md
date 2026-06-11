# TokenHub 用户角色与权限系统

## 角色定义

| 角色 | 值 | 说明 | 权限范围 |
|------|-----|------|----------|
| 游客 | 0 | 无需注册 | 仅可访问公开 API |
| 普通用户 | 1 | 注册用户 | 基本功能：管理自己的 API 密钥、查看使用统计 |
| 管理员 | 10 | 系统管理员 | 可管理普通用户、配置权益分组 |
| 超级管理员 | 100 | 系统唯一 | 可管理所有用户、配置系统设置 |

## 角色层级关系

```
超级管理员 (100)
    └── 管理员 (10)
        └── 普通用户 (1)
            └── 游客 (0)
```

## 权限矩阵

| 功能 | 游客 | 普通用户 | 管理员 | 超级管理员 |
|------|------|----------|--------|------------|
| 访问公开 API | ✅ | ✅ | ✅ | ✅ |
| 注册账户 | ✅ | - | - | - |
| 登录 | ✅ | ✅ | ✅ | ✅ |
| 管理自己的 API 密钥 | - | ✅ | ✅ | ✅ |
| 查看使用统计 | - | ✅ | ✅ | ✅ |
| 修改个人资料 | - | ✅ | ✅ | ✅ |
| 启用 2FA | - | ✅ | ✅ | ✅ |
| 管理普通用户 | - | - | ✅ | ✅ |
| 配置权益分组 | - | - | ✅ | ✅ |
| 管理管理员 | - | - | - | ✅ |
| 系统设置 | - | - | - | ✅ |

## 系统初始化

### 创建超级管理员

首次部署时，需要通过 API 创建超级管理员：

```bash
curl -X POST http://localhost:3000/api/setup \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'
```

**注意**：
- 超级管理员只能创建一次
- 创建后无法通过 API 再次创建
- 超级管理员可以自行修改密码

### 检查是否已初始化

```bash
curl http://localhost:3000/api/setup
```

返回：
```json
{
  "success": true,
  "setup": true  // true 表示已创建超级管理员
}
```

## 用户管理 API

### 管理员操作

| 操作 | API | 说明 |
|------|-----|------|
| 获取用户列表 | `GET /api/admin/user` | 分页获取所有用户 |
| 搜索用户 | `GET /api/admin/user/search?keyword=xxx` | 按关键词搜索 |
| 创建用户 | `POST /api/admin/user` | 管理员创建普通用户 |
| 编辑用户 | `PUT /api/admin/user/:id` | 修改用户信息 |
| 删除用户 | `DELETE /api/admin/user/:id` | 永久删除用户 |
| 禁用用户 | `POST /api/admin/user/manage` | action: "disable" |
| 启用用户 | `POST /api/admin/user/manage` | action: "enable" |
| 提升为管理员 | `POST /api/admin/user/manage` | action: "promote" (仅超级管理员) |
| 降级为普通用户 | `POST /api/admin/user/manage` | action: "demote" |
| 调整额度 | `POST /api/admin/user/manage` | action: "add_quota" |

### 用户自助操作

| 操作 | API | 说明 |
|------|-----|------|
| 获取个人信息 | `GET /api/user/self` | 获取当前用户信息 |
| 修改个人资料 | `PUT /api/user/self` | 修改显示名称、邮箱 |
| 修改密码 | `PUT /api/user/self` | 需要提供原密码 |
| 删除账户 | `DELETE /api/user/self` | 软删除自己的账户 |

## 权益分组

权益分组用于控制用户的 API 访问权限和配额。

### 预设分组

| 分组 | 额度 | 频率限制 | 说明 |
|------|------|----------|------|
| default | 100,000 | 60次/分钟 | 默认分组 |
| vip | 可配置 | 可配置 | VIP 用户 |
| svip | 可配置 | 可配置 | SVIP 用户 |

### 管理 API

| 操作 | API | 说明 |
|------|-----|------|
| 获取所有分组 | `GET /api/admin/privilege-group` | 管理员可访问 |
| 创建分组 | `POST /api/admin/privilege-group` | 管理员可访问 |
| 编辑分组 | `PUT /api/admin/privilege-group/:id` | 管理员可访问 |
| 删除分组 | `DELETE /api/admin/privilege-group/:id` | 管理员可访问 |

## 2FA 双因素认证

### 启用流程

1. 调用 `POST /api/user/2fa/enable` 获取二维码
2. 使用 Google Authenticator 等应用扫描二维码
3. 输入 6 位验证码
4. 调用 `POST /api/user/2fa/verify` 验证并启用

### 登录流程

1. 输入用户名和密码
2. 如果启用了 2FA，系统返回 `require_2fa: true`
3. 输入 2FA 验证码
4. 调用 `POST /api/user/login/2fa` 完成登录

## Passkey (WebAuthn)

Passkey 支持使用硬件密钥、生物识别等进行无密码登录。

### 启用流程

1. 调用 `POST /api/user/passkey/enable`
2. 完成 WebAuthn 注册流程
3. 调用 `POST /api/user/passkey/verify` 保存凭证

### 登录流程

1. 输入用户名
2. 调用 `POST /api/user/passkey/login/begin`
3. 完成 WebAuthn 认证
4. 调用 `POST /api/user/passkey/login/complete` 完成登录

## OAuth 登录

支持的 OAuth 提供商：
- GitHub
- Discord
- Google (OIDC)

### 绑定流程

1. 调用 `GET /api/oauth/{provider}/login` 开始 OAuth 流程
2. 在提供商页面授权
3. 回调自动绑定到当前账户

## 数据库表

| 表名 | 说明 |
|------|------|
| users | 用户信息 |
| privilege_groups | 权益分组 |
| two_factor_auths | 2FA 配置 |
| passkey_credentials | Passkey 凭证 |
| o_auth_bindings | OAuth 绑定 |
