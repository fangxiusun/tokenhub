# oauth.go (controller) 代码阅读文档

## 1. 全局总结
oauth.go 是 OAuth 统一控制器，处理所有 OAuth 提供商的登录和绑定逻辑。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `i18n` - 国际化
- `model` - 数据模型
- `oauth` - OAuth 提供商
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/api-router.go` - 管理 API 路由

## 3. 类型定义
### 3.1 函数
- `GenerateOAuthCode` - 生成 OAuth 授权码
- `HandleOAuth` - 处理 OAuth 回调

## 4. 函数详解
### 4.1 HandleOAuth
- **签名**: `func HandleOAuth(c *gin.Context)`
- **职责**: 统一 OAuth 回调处理
- **逻辑流程**:
  1. 从查询参数获取授权码和状态
  2. 验证状态参数
  3. 调用提供商 ExchangeToken 交换令牌
  4. 获取用户信息
  5. 查找或创建用户
  6. 设置会话

## 5. 关键逻辑分析
- **统一入口**: 所有 OAuth 提供商使用相同的回调处理逻辑
- **状态验证**: 防止 CSRF 攻击
- **用户创建**: 首次登录自动创建用户

## 6. 关联文件
- `oauth/` - OAuth 提供商实现
- `model/user.go` - 用户模型
