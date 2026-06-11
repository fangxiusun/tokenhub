# rate_limit.go 代码阅读文档

## 1. 全局总结

该文件实现模型请求速率限制配置，支持按用户组设置不同的请求次数限制。

## 2. 依赖关系

- `encoding/json` — JSON 序列化/反序列化
- `fmt` — 错误格式化
- `math` — 数学常量
- `sync` — 互斥锁
- `github.com/QuantumNous/new-api/common` — 系统日志

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelRequestRateLimitEnabled` | `bool` | 是否启用模型请求速率限制 |
| `ModelRequestRateLimitDurationMinutes` | `int` | 限制时间窗口（分钟） |
| `ModelRequestRateLimitCount` | `int` | 全局请求次数限制 |
| `ModelRequestRateLimitSuccessCount` | `int` | 全局成功请求次数限制 |
| `ModelRequestRateLimitGroup` | `map[string][2]int` | 按用户组的限制，值为 [总次数, 成功次数] |
| `ModelRequestRateLimitMutex` | `sync.RWMutex` | 读写锁 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `ModelRequestRateLimitGroup2JSONString` | `func ModelRequestRateLimitGroup2JSONString() string` | 将用户组限制序列化为 JSON |
| `UpdateModelRequestRateLimitGroupByJSONString` | `func UpdateModelRequestRateLimitGroupByJSONString(jsonStr string) error` | 从 JSON 更新用户组限制 |
| `GetGroupRateLimit` | `func GetGroupRateLimit(group string) (totalCount, successCount int, found bool)` | 获取指定用户组的限制 |
| `CheckModelRequestRateLimitGroup` | `func CheckModelRequestRateLimitGroup(jsonStr string) error` | 校验速率限制配置的合法性 |

## 5. 关键逻辑分析

- 使用 `sync.RWMutex` 保护 `ModelRequestRateLimitGroup` 的并发访问
- `CheckModelRequestRateLimitGroup` 校验值不能为负数且不超过 `math.MaxInt32`
- `GetGroupRateLimit` 返回 [总次数, 成功次数, 是否找到]

## 6. 关联文件

- `middleware/rate-limit.go` — 速率限制中间件
- `common/rate-limit.go` — 速率限制工具函数
