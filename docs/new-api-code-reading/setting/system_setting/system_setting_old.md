# system_setting_old.go 代码阅读文档

## 1. 全局总结

该文件定义旧版系统设置变量，包括服务器地址和 Worker 配置。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `ServerAddress` | `string` | `"http://localhost:3000"` | 服务器地址 |
| `WorkerUrl` | `string` | `""` | Worker URL |
| `WorkerValidKey` | `string` | `""` | Worker 验证密钥 |
| `WorkerAllowHttpImageRequestEnabled` | `bool` | `false` | 是否允许 Worker HTTP 图片请求 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `EnableWorker` | `func EnableWorker() bool` | 检查 Worker 是否启用 |

## 5. 关键逻辑分析

- `EnableWorker` 通过检查 `WorkerUrl` 是否非空来判断
- 该文件保留向后兼容

## 6. 关联文件

- `setting/system_setting/passkey.go` — 使用 ServerAddress
- `middleware/auth.go` — 使用 ServerAddress
