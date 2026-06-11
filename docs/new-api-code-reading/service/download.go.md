# download.go 代码阅读文档

## 1. 全局总结

该文件提供文件下载功能，支持通过 Worker 代理或直接下载两种模式。包含 SSRF 防护验证，确保下载请求的安全性。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | URL 验证、日志 |
| `system_setting` | Worker 配置、获取设置 |

## 3. 类型定义

### `WorkerRequest`
Worker 请求结构体：`URL`、`Key`、`Method`、`Headers`、`Body`

## 4. 函数详解

### `DoWorkerRequest(req) (*http.Response, error)`
- 通过 Worker 代理发送请求
- 验证 Worker 是否启用
- 非 HTTPS 请求需要额外配置允许
- SSRF 防护验证

### `DoDownloadRequest(originUrl, reason...) (resp, err)`
- 统一下载入口
- Worker 模式：通过 Worker 代理
- 直接模式：SSRF 验证后直接请求

## 5. 关键逻辑分析

1. **双模式支持**：Worker 代理 / 直接下载
2. **SSRF 防护**：两种模式都进行 URL 安全验证
3. **HTTPS 限制**：Worker 默认仅支持 HTTPS

## 6. 关联文件

- `http_client.go` — HTTP 客户端
- `file_service.go` — 文件服务
