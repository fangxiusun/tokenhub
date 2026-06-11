# console_migrate.go 代码阅读文档

## 1. 全局总结

该文件用于将旧版控制台配置迁移到新的 `console_setting.*` 命名空间。处理 ApiInfo、Announcements、FAQ 和 Uptime Kuma 等配置项的迁移。标记为即将删除的临时文件。

## 2. 依赖关系

- `common` — JSON 操作、日志
- `model` — Option 模型操作
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `MigrateConsoleSetting(c *gin.Context)`
执行控制台设置迁移。流程：读取全部 Option → 逐个处理旧键（ApiInfo、Announcements、FAQ、UptimeKumaUrl/Slug）→ 写入新键 → 删除旧键 → 重新加载 OptionMap。

## 5. 关键逻辑分析

- ApiInfo 数组限制最多 50 条
- FAQ 字段名标准化（`question`/`answer`）
- Uptime Kuma 迁移到 groups 结构（包含 id、categoryName、url、slug、description）
- 迁移完成后删除旧键记录

## 6. 关联文件

- `model/option.go` — Option 模型
- `setting/console_setting/` — 新版控制台设置
