# http.go 代码阅读文档

## 1. 全局总结

该文件提供 HTTP 响应处理的辅助函数，包括响应体关闭、上游响应头复制、以及字节流复制到客户端。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 请求 ID 键名 |
| `logger` | 错误日志 |
| `gin` | HTTP 上下文 |

## 3. 函数详解

### `CloseResponseBodyGracefully(httpResponse)`
安全关闭响应体，忽略错误

### `ShouldCopyUpstreamHeader(c, k, v) bool`
判断是否应该复制上游响应头：
- `Content-Length`：不复制（单独管理）
- `X-Oneapi-Request-Id`：不复制，保存到 context

### `IOCopyBytesGracefully(c, src, data)`
将字节数据复制到客户端响应：
1. 复制上游响应头
2. 手动设置 Content-Length
3. 写入状态码
4. 复制响应体

## 4. 关键逻辑分析

1. **延迟设头**：先解析响应体再设头，避免错误时头已发送
2. **请求 ID 保护**：保留本地实例的请求 ID

## 5. 关联文件

- `relay/` — 上游响应处理
