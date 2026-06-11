# request_body_limit.go 代码阅读文档

## 1. 全局总结

本文件提供匿名用户请求体大小限制的功能。通过读取配置常量，将 KB 单位的限制转换为字节数返回。设计简洁，仅包含一个常量和一个函数，用于控制匿名用户的请求体大小上限。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/constant` | 提供 `AnonymousRequestBodyLimitKB` 配置常量 |

**项目内依赖：**
- `constant.AnonymousRequestBodyLimitKB` — 匿名用户请求体限制（KB）

## 3. 类型定义

### 常量

```go
const defaultAnonymousRequestBodyLimitKB = 512
```

默认匿名用户请求体限制为 512 KB。当配置值小于 0 时使用此默认值。

## 4. 函数详解

### GetAnonymousRequestBodyLimitBytes() int64

```go
func GetAnonymousRequestBodyLimitBytes() int64
```

获取匿名用户请求体大小限制（字节单位）。

**逻辑流程：**
1. 读取 `constant.AnonymousRequestBodyLimitKB` 配置值
2. 如果配置值 < 0，使用默认值 512 KB
3. 将 KB 值左移 10 位（乘以 1024）转换为字节数
4. 返回 int64 类型的字节数

**返回值：** 请求体大小限制，单位为字节。

**示例：**
- 配置值为 512 → 返回 `512 * 1024 = 524288` 字节
- 配置值为 -1 → 使用默认值 512，返回 `524288` 字节
- 配置值为 1024 → 返回 `1048576` 字节（1 MB）

## 5. 关键逻辑分析

### 5.1 位运算转换

使用 `<< 10` 位运算代替乘法进行 KB 到 Bytes 的转换：
- `int64(limitKB) << 10` 等价于 `int64(limitKB) * 1024`
- 位运算效率更高，且语义清晰

### 5.2 默认值策略

当配置值小于 0 时回退到默认值（512 KB）：
- 这是一种防御性编程，防止配置错误导致限制过大或过小
- 默认值 512 KB 对于大多数 API 请求来说是合理的上限

### 5.3 设计意图

本函数用于保护系统免受超大请求体的攻击或误用：
- 匿名用户（未登录用户）的请求体大小受到更严格的限制
- 通过配置常量可灵活调整限制大小

## 6. 关联文件

- `new-api/constant/` — 定义 `AnonymousRequestBodyLimitKB` 常量
- `new-api/middleware/` — 可能使用此函数进行请求体大小校验
- `new-api/controller/` — 请求处理器，可能调用此函数验证请求
