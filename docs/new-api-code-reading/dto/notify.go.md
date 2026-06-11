# notify.go 代码阅读文档

## 1. 全局摘要

该文件定义了通知相关的数据结构和常量。包含通知结构体 `Notify`、通知类型常量，以及通知构造函数。主要用于系统内部通知，如配额超限、渠道更新、渠道测试等场景。

## 2. 依赖

无外部包依赖。

## 3. 类型定义

### Notify 结构体
通知数据结构：
- `Type` (string)：通知类型
- `Title` (string)：通知标题
- `Content` (string)：通知内容
- `Values` ([]interface{})：模板值数组

### 常量

**ContentValueParam**：内容值参数模板 `"{{value}}"`。

**通知类型常量**：
- `NotifyTypeQuotaExceed`："quota_exceed" - 配额超限通知
- `NotifyTypeChannelUpdate`："channel_update" - 渠道更新通知
- `NotifyTypeChannelTest`："channel_test" - 渠道测试通知

## 4. 函数详情

### NewNotify()
```go
func NewNotify(t string, title string, content string, values []interface{}) Notify
```
**功能**：创建通知实例。

**参数**：
- `t`：通知类型
- `title`：通知标题
- `content`：通知内容（支持 `{{value}}` 模板）
- `values`：模板值数组

**返回**：初始化后的 `Notify` 结构体。

## 5. 关键逻辑分析

1. **模板机制**：通知内容支持 `{{value}}` 占位符，通过 `Values` 数组提供实际值。

2. **类型分类**：预定义三种通知类型，覆盖配额、渠道更新、渠道测试场景。

3. **工厂函数**：`NewNotify` 提供便捷的通知创建方式，避免字段遗漏。

## 6. 相关文件

- `middleware/notify.go`：通知中间件，使用这些结构
- `controller/notify.go`：通知控制器
- `service/notify.go`：通知服务逻辑