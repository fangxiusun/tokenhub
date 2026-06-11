# user_settings.go 代码阅读文档

## 1. 全局摘要

该文件定义了用户设置的数据结构 `UserSetting`，包含通知配置、Webhook 设置、语言偏好等用户个人配置选项。

## 2. 依赖

无外部包依赖。

## 3. 类型定义

### UserSetting 结构体
用户设置结构：
- `NotifyType` (string)：额度预警类型
- `QuotaWarningThreshold` (float64)：额度预警阈值
- `WebhookUrl` (string)：Webhook 地址
- `WebhookSecret` (string)：Webhook 密钥
- `NotificationEmail` (string)：通知邮箱地址
- `BarkUrl` (string)：Bark 推送 URL
- `GotifyUrl` (string)：Gotify 服务器地址
- `GotifyToken` (string)：Gotify 应用令牌
- `GotifyPriority` (int)：Gotify 消息优先级
- `UpstreamModelUpdateNotifyEnabled` (bool)：是否接收上游模型更新通知
- `AcceptUnsetRatioModel` (bool)：是否接受未设置价格的模型
- `RecordIpLog` (bool)：是否记录 IP 日志
- `SidebarModules` (string)：左侧边栏模块配置
- `BillingPreference` (string)：扣费策略（订阅/钱包）
- `Language` (string)：用户语言偏好

### 通知类型常量
- `NotifyTypeEmail`："email" - 邮件通知
- `NotifyTypeWebhook`："webhook" - Webhook 通知
- `NotifyTypeBark`："bark" - Bark 推送
- `NotifyTypeGotify`："gotify" - Gotify 推送

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **多渠道通知**：支持邮件、Webhook、Bark、Gotify 等多种通知方式。

2. **安全配置**：`WebhookSecret` 用于 Webhook 验证，`RecordIpLog` 控制日志记录。

3. **个性化设置**：支持语言偏好、侧边栏配置等个性化选项。

4. **计费策略**：`BillingPreference` 支持订阅和钱包两种扣费模式。

## 6. 相关文件

- `model/user.go`：用户数据模型
- `controller/user.go`：用户控制器
- `service/notification.go`：通知服务