# codex_wham_usage.go 代码阅读文档

## 1. 全局总结

该文件提供了查询 Codex WHAM（Web Hosted AI Model）使用量的 API 调用功能。用于获取指定账户的 AI 模型使用统计信息。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `context` | 请求超时控制 |
| `net/http` | HTTP 客户端 |
| `io` | 响应体读取 |

## 3. 类型定义

无自定义类型。

## 4. 函数详解

### `FetchCodexWhamUsage(ctx, client, baseURL, accessToken, accountID) (int, []byte, error)`
- 向 `/backend-api/wham/usage` 发送 GET 请求
- 设置 `Authorization: Bearer {token}` 和 `chatgpt-account-id` 头
- 设置 `originator: codex_cli_rs` 头
- 返回状态码和响应体

## 5. 关键逻辑分析

1. **参数校验**：严格检查所有必要参数
2. **URL 规范化**：移除 baseURL 末尾斜杠
3. **原始响应**：返回原始字节，由调用方解析

## 6. 关联文件

- `codex_oauth.go` — OAuth 令牌获取
- `codex_credential_refresh.go` — 凭证刷新
