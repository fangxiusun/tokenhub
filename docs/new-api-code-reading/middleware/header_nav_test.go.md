# header_nav_test.go 代码阅读文档

## 1. 全局总结

该文件是 `header_nav.go` 的单元测试文件，使用 Go 标准测试框架和 `testify/require` 断言库，测试 `HeaderNavModuleAuth` 和 `HeaderNavModulePublicOrUserAuth` 两个中间件在各种配置场景下的行为。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `net/http` | HTTP 状态码和请求方法 |
| `net/http/httptest` | 测试用 HTTP 服务器 |
| `testing` | Go 测试框架 |
| `github.com/QuantumNous/new-api/common` | 配置存储、常量 |
| `github.com/gin-contrib/sessions` | Session 管理 |
| `github.com/gin-contrib/sessions/cookie` | Cookie-based Session 存储 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |
| `github.com/stretchr/testify/require` | 测试断言 |

## 3. 辅助函数

### `withHeaderNavModules(t *testing.T, raw string)`

- **功能**：设置测试用的 `HeaderNavModules` 配置，并在测试结束后恢复原始值。
- **实现**：通过 `t.Cleanup` 确保配置恢复，避免测试间干扰。

### `performHeaderNavRequest(t *testing.T, handler gin.HandlerFunc, authenticated bool) *httptest.ResponseRecorder`

- **功能**：执行一个带/不带认证的 HTTP 请求，返回响应记录器。
- **逻辑**：
  1. 创建 Gin 测试路由（含 session 中间件）
  2. 如果需要认证，先调用 `/login` 获取 session cookie
  3. 发送 `GET /api/test` 请求，携带认证信息

## 4. 测试用例

### `HeaderNavModuleAuth` 测试

| 测试名 | 场景 | 预期结果 |
|--------|------|----------|
| `TestHeaderNavModuleAuthAllowsDefaultPublicAccess` | 空配置（默认公开） | 200 OK |
| `TestHeaderNavModuleAuthRejectsDisabledPricing` | pricing 模块禁用 | 403 Forbidden |
| `TestHeaderNavModuleAuthRequiresLoginForPricing` | pricing 需要认证 | 401 Unauthorized |
| `TestHeaderNavModuleAuthRequiresLoginForRankings` | rankings 需要认证 | 401 Unauthorized |
| `TestHeaderNavModuleAuthRejectsLegacyDisabledModule` | rankings 旧格式禁用（`false`） | 403 Forbidden |

### `HeaderNavModulePublicOrUserAuth` 测试

| 测试名 | 场景 | 预期结果 |
|--------|------|----------|
| `TestHeaderNavModulePublicOrUserAuthAllowsDefaultPublicAccess` | 空配置（默认公开） | 200 OK |
| `TestHeaderNavModulePublicOrUserAuthRequiresLoginWhenDisabled` | 禁用时未登录 | 401 Unauthorized |
| `TestHeaderNavModulePublicOrUserAuthAllowsLoggedInWhenDisabled` | 禁用时已登录 | 200 OK |
| `TestHeaderNavModulePublicOrUserAuthRequiresLoginWhenRequireAuth` | 需要认证时未登录 | 401 Unauthorized |
| `TestHeaderNavModulePublicOrUserAuthRequiresLoginForLegacyDisabledModule` | 旧格式禁用时未登录 | 401 Unauthorized |

## 5. 关键逻辑分析

- **测试隔离**：使用 `withHeaderNavModules` 的 `t.Cleanup` 机制确保每个测试结束后恢复原始配置。
- **Session 模拟**：使用 cookie-based session store 模拟用户登录状态。
- **认证模拟**：通过设置 `New-Api-User` header 和 session cookie 模拟已认证请求。
- **覆盖全面**：测试了空配置、完整对象配置、旧格式布尔值配置等多种场景。

## 6. 关联文件

- `middleware/header_nav.go` — 被测试的源文件
- `common/option.go` — `OptionMap` 配置存储
