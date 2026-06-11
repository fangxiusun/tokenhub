# header_nav.go 代码阅读文档

## 1. 全局总结

该文件实现了头部导航模块的访问控制中间件，允许管理员通过配置动态控制各前端模块（如 pricing、rankings 等）的启用状态和认证要求。支持两种中间件模式：`HeaderNavModuleAuth`（模块级认证控制）和 `HeaderNavModulePublicOrUserAuth`（公开或认证访问）。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化错误消息 |
| `net/http` | HTTP 状态码 |
| `strings` | 字符串处理 |
| `github.com/QuantumNous/new-api/common` | 配置存储、JSON 解析、认证中间件 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

### `headerNavAccess`

```go
type headerNavAccess struct {
    Enabled     bool  // 模块是否启用
    RequireAuth bool  // 是否需要认证
}
```

## 4. 函数详解

### `getHeaderNavAccess(module string) headerNavAccess`

- **功能**：从 `OptionMap["HeaderNavModules"]` 读取模块配置，返回访问控制信息。
- **返回**：默认值为 `{Enabled: true, RequireAuth: false}`（公开访问）。

### `parseHeaderNavAccess(raw any, fallback headerNavAccess) headerNavAccess`

- **功能**：解析配置值，支持多种格式。
- **支持格式**：`bool`、`string`、`float64`、`map[string]any`。

### `parseHeaderNavBool(value any, fallback bool) bool`

- **功能**：将各种类型转换为布尔值，支持 `"true"/"1"` 和 `"false"/"0"` 字符串。

### `HeaderNavModuleAuth(module string) gin.HandlerFunc`

- **功能**：模块级认证控制中间件。
- **逻辑**：
  - 模块禁用 → 返回 403
  - 模块启用且 requireAuth=true → 强制认证（`UserAuth`）
  - 模块启用且 requireAuth=false → 尝试认证（`TryUserAuth`）

### `HeaderNavModulePublicOrUserAuth(module string) gin.HandlerFunc`

- **功能**：公开或认证访问中间件。
- **逻辑**：
  - 模块禁用或 requireAuth=true → 强制认证
  - 模块启用且 requireAuth=false → 尝试认证

## 5. 关键逻辑分析

- **配置来源**：从 `common.OptionMap` 读取 `HeaderNavModules` JSON 配置，使用读写锁保证并发安全。
- **灵活配置格式**：支持布尔值简写（`{"pricing": false}`）和完整对象（`{"pricing": {"enabled": true, "requireAuth": true}}`）。
- **中间件差异**：`HeaderNavModuleAuth` 在模块禁用时返回 403；`HeaderNavModulePublicOrUserAuth` 在模块禁用时强制认证（已登录用户仍可访问）。

## 6. 关联文件

- `middleware/auth.go` — `UserAuth` 和 `TryUserAuth` 认证中间件实现
- `common/option.go` — `OptionMap` 配置存储
