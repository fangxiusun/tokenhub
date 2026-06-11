# setup.go 代码阅读文档

## 1. 全局总结

该文件实现了系统初始化（Setup）功能。首次部署时创建 Root 管理员账号，配置运营模式。

## 2. 依赖关系

- `common` — 密码哈希、数据库类型
- `constant` — Setup 标志
- `model` — 用户模型、Setup 模型
- `setting/operation_setting` — 运营模式
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `Setup` | 系统状态响应 |
| `SetupRequest` | 初始化请求 |

## 4. 函数详解

### `GetSetup(c *gin.Context)`
获取系统初始化状态。返回是否已完成初始化、Root 用户是否存在、数据库类型。

### `PostSetup(c *gin.Context)`
执行系统初始化。创建 Root 用户（角色为 root，初始额度 1 亿），设置运营模式（自用模式、演示站点），记录初始化时间。

## 5. 关键逻辑分析

- 仅在系统未初始化时允许执行
- Root 用户名最长 12 字符，密码最少 8 字符
- Root 用户初始额度：100000000
- 初始化完成后设置 `constant.Setup = true`

## 6. 关联文件

- `model/user.go` — 用户模型
- `model/setup.go` — Setup 模型
