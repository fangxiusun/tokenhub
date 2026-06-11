# auth.go 代码阅读文档

## 1. 全局总结
auth.go 是中间件层的核心鉴权文件，实现了多种鉴权机制：Session 鉴权、Access Token 鉴权、API Token 鉴权。提供了 UserAuth、AdminAuth、RootAuth、TokenAuth 四个中间件函数，分别用于不同权限级别的接口保护。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具（JSON、日志、验证）
- `constant` - 常量定义（上下文键、角色常量）
- `i18n` - 国际化翻译
- `logger` - 结构化日志
- `model` - 数据模型（用户、Token 验证）
- `service` - 业务逻辑（渠道路由）
- `setting/ratio_setting` - 比例设置
- `types` - 类型定义
- `gin-contrib/sessions` - Session 管理
- `gin-gonic/gin` - Web 框架
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `router/api-router.go` - 管理 API 路由
- `router/relay-router.go` - 中转路由
- `controller/` - 所有控制器

## 3. 类型定义
### 3.1 函数
- `validUserInfo(username string, role int) bool` - 验证用户信息有效性
- `authHelper(c *gin.Context, minRole int)` - 核心鉴权辅助函数
- `UserAuth() gin.HandlerFunc` - 普通用户鉴权中间件
- `AdminAuth() gin.HandlerFunc` - 管理员鉴权中间件
- `RootAuth() gin.HandlerFunc` - 超级管理员鉴权中间件
- `TokenAuth() gin.HandlerFunc` - API Token 鉴权中间件

## 4. 函数详解
### 4.1 authHelper
- **签名**: `func authHelper(c *gin.Context, minRole int)`
- **职责**: 核心鉴权逻辑，支持 Session 和 Access Token 两种方式
- **逻辑流程**:
  1. 尝试从 Session 获取用户信息
  2. 如果 Session 无用户，检查 Authorization 头的 Access Token
  3. 验证用户状态和角色
  4. 检查 New-Api-User 头（用户切换）
  5. 设置上下文值供后续使用

### 4.2 TokenAuth
- **签名**: `func TokenAuth() gin.HandlerFunc`
- **职责**: API Token 鉴权，用于中转接口
- **逻辑流程**:
  1. 从多个来源提取 Token（Authorization 头、x-api-key 头、查询参数）
  2. 验证 Token 状态、过期时间、剩余配额
  3. 检查 IP 限制
  4. 检查模型限制
  5. 设置 Token 上下文信息

## 5. 关键逻辑分析
- **多源 Token 提取**: 支持 Authorization 头、x-api-key 头、Sec-WebSocket-Protocol 头、query 参数
- **用户切换**: 通过 New-Api-User 头支持管理员切换到其他用户
- **信任模式**: 高余额用户可跳过预扣费
- **错误处理**: 数据库错误返回 500，其他错误返回 200 + success:false

## 6. 关联文件
- `model/user.go` - 用户模型验证
- `model/token.go` - Token 验证
- `constant/context_key.go` - 上下文键定义
- `i18n/keys.go` - 翻译键定义
