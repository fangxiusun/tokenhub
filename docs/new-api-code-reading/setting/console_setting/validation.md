# validation.go 代码阅读文档

## 1. 全局总结

该文件实现控制台设置的输入校验，包括 API 信息、系统公告、FAQ、Uptime Kuma 分组的格式和内容验证，以及数据获取和排序功能。

## 2. 依赖关系

- `encoding/json` — JSON 解析
- `fmt` — 错误格式化
- `net/url` — URL 解析
- `regexp` — 正则表达式
- `sort` — 排序
- `strings` — 字符串处理
- `time` — 时间处理

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `urlRegex` | `*regexp.Regexp` | URL 格式正则 |
| `dangerousChars` | `[]string` | 危险内容检测关键词列表 |
| `validColors` | `map[string]bool` | 合法颜色值集合 |
| `slugRegex` | `*regexp.Regexp` | Slug 格式正则 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `ValidateConsoleSettings` | `func ValidateConsoleSettings(settingsStr string, settingType string) error` | 统一校验入口 |
| `validateApiInfo` | `func validateApiInfo(apiInfoStr string) error` | 校验 API 信息 |
| `validateAnnouncements` | `func validateAnnouncements(announcementsStr string) error` | 校验系统公告 |
| `validateFAQ` | `func validateFAQ(faqStr string) error` | 校验 FAQ |
| `validateUptimeKumaGroups` | `func validateUptimeKumaGroups(groupsStr string) error` | 校验 Uptime Kuma 分组 |
| `GetApiInfo` | `func GetApiInfo() []map[string]interface{}` | 获取 API 信息列表 |
| `GetAnnouncements` | `func GetAnnouncements() []map[string]interface{}` | 获取公告列表（按时间倒序） |
| `GetFAQ` | `func GetFAQ() []map[string]interface{}` | 获取 FAQ 列表 |
| `GetUptimeKumaGroups` | `func GetUptimeKumaGroups() []map[string]interface{}` | 获取 Uptime Kuma 分组列表 |

## 5. 关键逻辑分析

- API 信息校验：URL 格式、必填字段、长度限制、颜色合法性、XSS 防护
- 公告校验：数量上限 100、必填 content 和 publishDate、日期格式 RFC3339、类型合法性
- FAQ 校验：数量上限 100、question/answer 必填、长度限制
- Uptime Kuma 校验：数量上限 20、名称唯一性、URL 格式、Slug 格式
- 危险内容检测：`<script`、`<iframe`、`javascript:`、`onload=`、`onerror=`、`onclick=`
- `GetAnnouncements` 使用 `sort.SliceStable` 按发布时间倒序排列

## 6. 关联文件

- `setting/console_setting/config.go` — 控制台配置定义
- `controller/console.go` — 控制台接口
