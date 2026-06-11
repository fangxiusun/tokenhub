# theme.go 代码阅读文档

## 1. 全局总结

该文件定义前端主题设置，并提供与 `common` 包的同步机制。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/common` — 主题设置同步
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `ThemeSettings` | `Frontend` | `string` | `"classic"` | 前端主题 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetThemeSettings` | `func GetThemeSettings() *ThemeSettings` | 获取主题设置 |
| `UpdateAndSyncTheme` | `func UpdateAndSyncTheme()` | 更新并同步主题到 common 包 |

## 5. 关键逻辑分析

- 默认主题为 "classic"
- `syncThemeToCommon` 在 `init()` 和配置更新后调用
- 主题设置影响前端渲染

## 6. 关联文件

- `web/default/` — 默认前端主题
- `web/classic/` — 经典前端主题
- `common/theme.go` — 主题存储
