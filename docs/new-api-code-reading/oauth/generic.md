# generic.go 代码阅读文档

## 1. 全局概述

本文件实现了通用 OAuth 提供者（GenericOAuthProvider），支持通过管理员配置的自定义 OAuth 提供者登录。提供灵活的字段映射、访问策略评估和多种认证风格支持。

## 2. 依赖关系

- `context` — 上下文
- `encoding/base64` — Base64 编码
- `encoding/json` — JSON 编码（别名 `stdjson`）
- `errors` — 错误处理
- `fmt` — 格式化输出
- `io` — I/O 操作
- `net/http` — HTTP 客户端
- `net/url` — URL 编码
- `regexp` — 正则表达式
- `strconv` — 字符串转换
- `strings` — 字符串操作
- `time` — 超时控制
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/QuantumNous/new-api/i18n` — 国际化
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/model` — 数据模型
- `github.com/QuantumNous/new-api/setting/system_setting` — 系统设置
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/samber/lo` — Lodash 风格工具
- `github.com/tidwall/gjson` — JSON 路径查询

## 3. 类型定义

### AuthStyle 常量

```go
const (
    AuthStyleAutoDetect = 0
    AuthStyleInParams   = 1
    AuthStyleInHeader   = 2
)
```

### GenericOAuthProvider 结构体

```go
type GenericOAuthProvider struct {
    config *model.CustomOAuthProvider
}
```

### accessPolicy 结构体

```go
type accessPolicy struct {
    Logic      string            `json:"logic"`
    Conditions []accessCondition `json:"conditions"`
    Groups     []accessPolicy    `json:"groups"`
}
```

### accessCondition 结构体

```go
type accessCondition struct {
    Field string `json:"field"`
    Op    string `json:"op"`
    Value any    `json:"value"`
}
```

### accessPolicyFailure 结构体

```go
type accessPolicyFailure struct {
    Field    string
    Op       string
    Expected any
    Current  any
}
```

## 4. 函数详情

### 构造函数

- `NewGenericOAuthProvider(config *model.CustomOAuthProvider) *GenericOAuthProvider`

### Provider 接口实现

- `GetName() string` — 返回配置的名称
- `IsEnabled() bool` — 返回配置的启用状态
- `GetConfig() *model.CustomOAuthProvider` — 获取配置
- `ExchangeToken(...)` — 令牌交换
- `GetUserInfo(...)` — 获取用户信息
- `IsUserIDTaken(...)` — 检查用户 ID 是否已占用
- `FillUserByProviderID(...)` — 通过 provider 用户 ID 填充用户
- `SetProviderUserID(...)` — 设置 provider 用户 ID
- `GetProviderPrefix() string` — 返回 `{slug}_` 前缀
- `GetProviderId() int` — 获取 provider ID
- `IsGenericProvider() bool` — 返回 true

### 访问策略函数

- `parseAccessPolicy(raw string) (*accessPolicy, error)` — 解析策略 JSON
- `validateAccessPolicy(policy *accessPolicy) error` — 验证策略
- `validateAccessCondition(condition *accessCondition, index int) error` — 验证条件
- `evaluateAccessPolicy(body string, policy *accessPolicy) (bool, *accessPolicyFailure)` — 评估策略
- `evaluateAccessCondition(body string, cond accessCondition) (bool, *accessPolicyFailure)` — 评估条件

### 工具函数

- `normalizeAuthorizationTokenType(tokenType string) string` — 标准化 Token 类型
- `gjsonResultToValue(result gjson.Result) any` — gjson 结果转换
- `compareAny(left any, right any) int` — 通用比较
- `toFloat(v any) (float64, bool)` — 转换为浮点数
- `valueInSlice(current any, expected any) bool` — 检查值是否在列表中
- `containsValue(current any, expected any) bool` — 检查包含关系
- `renderAccessDeniedMessage(...)` — 渲染访问拒绝消息

## 5. 关键逻辑分析

### 认证风格

| 风格 | 常量 | 说明 |
|------|------|------|
| AutoDetect | 0 | 自动检测（默认使用 Params） |
| InParams | 1 | client_id/client_secret 作为 POST 参数 |
| InHeader | 2 | 作为 Basic Auth 头发送 |

### 字段映射

使用 `gjson` 库支持 JSONPath 风格的字段映射：
- `UserIdField` — 用户 ID 字段路径
- `UsernameField` — 用户名字段路径
- `DisplayNameField` — 显示名称字段路径
- `EmailField` — 邮箱字段路径

### 访问策略系统

#### 策略结构

```json
{
  "logic": "and",
  "conditions": [
    {"field": "trust_level", "op": "gte", "value": 2},
    {"field": "active", "op": "eq", "value": true}
  ],
  "groups": [
    {
      "logic": "or",
      "conditions": [...]
    }
  ]
}
```

#### 支持的操作符

| 操作符 | 说明 |
|--------|------|
| `eq` / `ne` | 等于 / 不等于 |
| `gt` / `gte` | 大于 / 大于等于 |
| `lt` / `lte` | 小于 / 小于等于 |
| `in` / `not_in` | 在列表中 / 不在列表中 |
| `contains` / `not_contains` | 包含 / 不包含 |
| `exists` / `not_exists` | 存在 / 不存在 |

#### 评估逻辑

- `and` 逻辑：所有条件和组都必须满足
- `or` 逻辑：任一条件或组满足即可
- 支持嵌套组（递归评估）

### 令牌交换

1. 构造 redirect_uri：`{ServerAddress}/oauth/{slug}`
2. 根据 AuthStyle 决定认证方式
3. 支持 JSON 和 URL-encoded 两种响应格式
4. 错误响应通过 `error` 和 `error_description` 字段识别

### 用户信息获取

1. 使用 `gjson` 按配置的字段路径提取数据
2. 支持数字和字符串类型的用户 ID
3. 评估访问策略，不满足时返回 `AccessDeniedError`

### 消息模板渲染

`renderAccessDeniedMessage` 支持以下模板变量：
- `{{provider}}` — 提供者名称
- `{{field}}` — 失败的字段
- `{{op}}` — 失败的操作符
- `{{required}}` — 期望值
- `{{current}}` — 当前值
- `{{current.path}}` — 从原始响应中提取的值
- `{{required.path}}` — 从期望条件中提取的值

## 6. 相关文件

- `oauth/provider.go` — Provider 接口
- `oauth/types.go` — OAuthToken、OAuthUser、AccessDeniedError 类型
- `oauth/registry.go` — 注册表
- `model/custom_oauth.go` — CustomOAuthProvider 数据模型
- `setting/system_setting/` — 系统设置
