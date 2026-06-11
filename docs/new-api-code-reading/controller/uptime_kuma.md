# uptime_kuma.go 代码阅读文档

## 1. 全局总结

该文件实现了 Uptime Kuma 状态监控的集成，通过代理请求 Uptime Kuma API 获取监控状态和心跳数据，避免前端直接访问外部服务的 CORS 问题。

## 2. 依赖关系

- `setting/console_setting` — Uptime Kuma 配置
- `gin-gonic/gin` — HTTP 框架
- `golang.org/x/sync/errgroup` — 并发控制

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `Monitor` | 监控项（名称、正常运行时间、状态、分组） |

## 4. 函数详解

### 代理 API
代理前端请求到 Uptime Kuma 服务器，获取状态页面和心跳数据。

## 5. 关键逻辑分析

- 请求超时：30 秒（整体）、10 秒（HTTP）
- 支持多分组的 Uptime Kuma 状态页面
- 使用 errgroup 进行并发请求

## 6. 关联文件

- `setting/console_setting/` — Uptime Kuma 配置
