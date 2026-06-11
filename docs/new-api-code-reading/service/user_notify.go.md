# user_notify.go 代码阅读文档

## 1. 全局总结

该文件实现用户通知系统，支持多种通知渠道：邮件、Webhook、Bark、Gotify。提供统一的通知发送接口。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 邮件发送、URL 验证 |
| `dto` | 通知数据结构 |
| `model` | 用户查询 |
| `system_setting` | Worker 配置 |
| `http` | HTTP 请求 |

## 3. 类型定义

无自定义类型。

## 4. 函数详解

### `NotifyRootUser(t, subject, content)`
通知管理员

### `NotifyUpstreamModelUpdateWatchers(subject, content)`
通知上游模型更新的订阅者

### `NotifyUser(userId, userEmail, userSetting, data) error`
统一通知入口：
1. 检查通知频率限制
2. 根据通知类型分发

### 渠道实现
- `sendEmailNotify` — 邮件通知
- `SendWebhookNotify` — Webhook 通知（带 HMAC 签名）
- `sendBarkNotify` — Bark 推送通知
- `sendGotifyNotify` — Gotify 推送通知

## 5. 关键逻辑分析

1. **频率限制**：每个通知类型独立限制
2. **Worker 支持**：所有渠道都支持通过 Worker 代理发送
3. **SSRF 防护**：非 Worker 模式下验证 URL 安全性
4. **模板变量**：使用 `{{value}}` 占位符

## 6. 关联文件

- `webhook.go` — Webhook 签名
- `notify-limit.go` — 频率限制
- `download.go` — Worker 请求
