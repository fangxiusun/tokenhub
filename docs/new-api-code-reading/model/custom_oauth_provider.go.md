# custom_oauth_provider.go 代码阅读文档

## 1. 全局总结
custom_oauth_provider.go 定义了自定义 OAuth 提供商模型，支持动态添加 OAuth 提供商。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `controller/custom_oauth.go` - 自定义 OAuth 控制器
- `oauth/generic.go` - 通用 OAuth 提供商

## 3. 类型定义
### 3.1 结构体
- `CustomOAuthProvider` - 自定义 OAuth 提供商结构体
  - `Id int` - 提供商 ID
  - `Name string` - 提供商名称
  - `DisplayName string` - 显示名称
  - `ClientId string` - 客户端 ID
  - `ClientSecret string` - 客户端密钥
  - `AuthorizationURL string` - 授权 URL
  - `TokenURL string` - 令牌 URL
  - `UserInfoURL string` - 用户信息 URL
  - `RedirectURL string` - 重定向 URL
  - `Scopes string` - 权限范围
  - `Status int` - 状态

## 4. 函数详解
### 4.1 GetAllCustomOAuthProviders
- **职责**: 获取所有自定义 OAuth 提供商

### 4.2 CreateCustomOAuthProvider
- **职责**: 创建自定义 OAuth 提供商

### 4.3 UpdateCustomOAuthProvider
- **职责**: 更新自定义 OAuth 提供商

### 4.4 DeleteCustomOAuthProvider
- **职责**: 删除自定义 OAuth 提供商

## 5. 关键逻辑分析
- **动态注册**: 运行时添加/删除 OAuth 提供商
- **配置存储**: OAuth 配置存储在数据库中

## 6. 关联文件
- `oauth/generic.go` - 通用 OAuth 实现
- `controller/custom_oauth.go` - 控制器
