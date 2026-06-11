# config.go 代码阅读文档

## 1. 全局总结

该文件定义控制台设置（ConsoleSetting）的数据结构和默认配置，管理 API 信息、Uptime Kuma、公告、FAQ 等控制台面板的配置。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `ConsoleSetting` | `ApiInfo` | `string` | API 信息（JSON 数组字符串） |
| | `UptimeKumaGroups` | `string` | Uptime Kuma 分组配置 |
| | `Announcements` | `string` | 系统公告 |
| | `FAQ` | `string` | 常见问题 |
| | `ApiInfoEnabled` | `bool` | 是否启用 API 信息面板 |
| | `UptimeKumaEnabled` | `bool` | 是否启用 Uptime Kuma 面板 |
| | `AnnouncementsEnabled` | `bool` | 是否启用系统公告面板 |
| | `FAQEnabled` | `bool` | 是否启用常见问答面板 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetConsoleSetting` | `func GetConsoleSetting() *ConsoleSetting` | 获取全局控制台配置实例 |

## 5. 关键逻辑分析

- 通过 `init()` 注册到 `config.GlobalConfig`，键名为 `"console_setting"`
- 所有配置项默认启用（`true`）
- 数据以 JSON 字符串形式存储，前端负责解析

## 6. 关联文件

- `setting/console_setting/validation.go` — 配置校验逻辑
- `controller/console.go` — 控制台相关接口
- `web/default/src/` — 前端控制台页面
