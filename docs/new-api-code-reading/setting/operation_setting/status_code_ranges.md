# status_code_ranges.go 代码阅读文档

## 1. 全局总结

该文件实现 HTTP 状态码范围管理，用于配置渠道自动禁用和自动重试的状态码规则。

## 2. 依赖关系

- `fmt` — 错误格式化
- `sort` — 排序
- `strconv` — 类型转换
- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/types` — 错误码类型

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `StatusCodeRange` | `Start` | `int` | 范围起始状态码 |
| | `End` | `int` | 范围结束状态码 |

| 变量名 | 说明 |
|--------|------|
| `AutomaticDisableStatusCodeRanges` | 自动禁用状态码范围（默认 401） |
| `AutomaticRetryStatusCodeRanges` | 自动重试状态码范围 |
| `alwaysSkipRetryStatusCodes` | 始终跳过重试的状态码（504, 524） |
| `alwaysSkipRetryCodes` | 始终跳过重试的错误码 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `ShouldDisableByStatusCode` | `func ShouldDisableByStatusCode(code int) bool` | 判断是否应禁用渠道 |
| `ShouldRetryByStatusCode` | `func ShouldRetryByStatusCode(code int) bool` | 判断是否应重试 |
| `IsAlwaysSkipRetryStatusCode` | `func IsAlwaysSkipRetryStatusCode(code int) bool` | 判断是否始终跳过重试 |
| `IsAlwaysSkipRetryCode` | `func IsAlwaysSkipRetryCode(errorCode types.ErrorCode) bool` | 判断错误码是否始终跳过重试 |
| `ParseHTTPStatusCodeRanges` | `func ParseHTTPStatusCodeRanges(input string) ([]StatusCodeRange, error)` | 解析状态码范围字符串 |

## 5. 关键逻辑分析

- 状态码范围格式：`401,403,500-599`
- 解析后自动排序并合并重叠范围
- 默认重试规则：1xx、3xx、4xx(除400/408)、5xx(除504/524)
- 504 和 524 始终跳过重试

## 6. 关联文件

- `relay/handler.go` — 使用重试规则
- `setting/operation_setting/operation_setting.go` — 运营设置
