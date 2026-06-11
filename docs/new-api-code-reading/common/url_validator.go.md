# url_validator.go 代码阅读文档

## 1. 全局总结
url_validator.go 实现了重定向 URL 安全验证，检查 URL 方案和域名是否在受信任列表中，防止开放重定向攻击。

## 2. 依赖关系
### 2.1 导入的包
- `net/url` - URL 解析
- `strings` - 字符串操作
- `common` - 公共工具

### 2.2 被引用的文件
- `controller/` - 控制器中验证重定向 URL

## 3. 类型定义
### 3.1 函数
- `ValidateRedirectURL(urlStr string) bool` - 验证重定向 URL
- `IsTrustedDomain(domain string) bool` - 检查域名是否受信任
- `GetTrustedRedirectDomains() []string` - 获取受信任域名列表

## 4. 函数详解
### 4.1 ValidateRedirectURL
- **签名**: `func ValidateRedirectURL(urlStr string) bool`
- **职责**: 验证 URL 是否安全可重定向
- **逻辑流程**:
  1. 解析 URL
  2. 检查方案（仅 http/https）
  3. 检查域名是否在受信任列表中
  4. 支持子域名匹配

## 5. 关键逻辑分析
- **子域名匹配**: 支持 `*.example.com` 格式的通配符
- **白名单机制**: 只允许受信任域名列表中的 URL

## 6. 关联文件
- `common/constants.go` - 受信任域名配置
- `controller/` - 使用此验证器
