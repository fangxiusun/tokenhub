# vendor_meta.go 代码阅读文档

## 1. 全局总结

`vendor_meta.go` 实现了供应商（Vendor）元数据的管理功能。供应商用于存储 AI 模型的上游厂商信息（如 OpenAI、Anthropic 等），供模型引用。提供基本的 CRUD 操作、名称唯一性检查和关键字搜索功能。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 时间戳获取（`GetTimestamp()`） |
| `gorm.io/gorm` | ORM 框架，提供软删除和查询能力 |

## 3. 类型定义

### Vendor（供应商）

```go
type Vendor struct {
    Id          int            // 主键
    Name        string         // 供应商名称（唯一索引，最大128字符）
    Description string         // 供应商描述
    Icon        string         // 图标名（@lobehub/icons 格式）
    Status      int            // 状态（1=启用，默认1）
    CreatedTime int64          // 创建时间
    UpdatedTime int64          // 更新时间
    DeletedAt   gorm.DeletedAt // 软删除标记
}
```

**唯一索引**：`(Name, DeletedAt)` — 保证未删除状态下供应商名称唯一。

## 4. 函数详解

### CRUD 操作

| 函数 | 说明 |
|------|------|
| `Insert()` | 创建供应商，自动设置创建时间和更新时间 |
| `Update()` | 更新供应商，自动刷新更新时间 |
| `Delete()` | 软删除供应商 |

### 查询

| 函数 | 说明 |
|------|------|
| `GetVendorByID(id)` | 按 ID 获取供应商 |
| `GetAllVendors(offset, limit)` | 分页获取所有供应商 |
| `SearchVendors(keyword, offset, limit)` | 按关键字搜索供应商（匹配名称和描述），返回分页结果和总数 |

### 检查

| 函数 | 说明 |
|------|------|
| `IsVendorNameDuplicated(id, name)` | 检查供应商名称是否重复（排除指定ID） |

## 5. 关键逻辑分析

### 名称唯一性
- 通过 `(Name, DeletedAt)` 复合唯一索引保证未删除状态下名称唯一
- `IsVendorNameDuplicated()` 在排除自身 ID 的情况下检查重复，用于编辑时的校验

### 图标规范
`Icon` 字段采用 `@lobehub/icons` 的图标命名规范，前端可直接渲染，无需额外映射。

### 搜索功能
`SearchVendors()` 使用 `LIKE` 模糊匹配名称和描述字段，支持分页和排序（按 ID 降序）。

### 软删除
使用 GORM 的 `DeletedAt` 字段实现软删除，删除后记录仍保留在数据库中，唯一索引会排除已删除记录。

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/channel.go` | 渠道模型可能引用 Vendor 信息 |
| `model/model.go` | 模型定义中关联供应商 |
| `controller/vendor.go` | 供应商管理 API 入口 |
| `web/` | 前端使用 `Icon` 字段渲染供应商图标 |
