# webhook.go 代码阅读文档

## 1. 全局总结

该文件实现 Webhook 通知发送功能，支持 HMAC-SHA256 签名验证。是通知系统的基础组件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | URL 验证 |
| `dto` | 通知数据结构 |
| `system_setting` | Worker 配置 |
| `crypto/hmac` | HMAC 签名 |
| `crypto/sha256` | SHA256 哈希 |

## 3. 类型定义

### `WebhookPayload`
Webhook 负载：Type、Title、Content、Values、Timestamp

## 4. 函数详解

### `generateSignature(secret, payload) string`
生成 HMAC-SHA256 签名

### `SendWebhookNotify(webhookURL, secret, data) error`
发送 Webhook 通知：
1. 处理模板变量
2. 构建负载
3. 序列化为 JSON
4. 生成签名
5. 发送请求（Worker 或直接）

## 5. 关键逻辑分析

1. **签名验证**：X-Webhook-Signature 头携带 HMAC-SHA256 签名
2. **Worker 支持**：支持通过 Worker 代理发送
3. **SSRF 防护**：非 Worker 模式下验证 URL

## 6. 关联文件

- `user_notify.go` — 通知发送
- `download.go` — Worker 请求
