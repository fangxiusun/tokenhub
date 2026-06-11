# page_info.go 代码阅读文档

## 1. 全局总结

本文件实现了分页查询的通用工具，提供了分页参数的结构体定义、索引计算方法，以及从 HTTP 请求中解析分页参数的函数 `GetPageQuery`。该函数支持多种查询参数名的兼容（`p`/`ps`/`size`/`page_size`），并对页码和页大小进行边界校验和默认值设置。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `strconv` | 将查询参数字符串转换为整数 |
| `github.com/gin-gonic/gin` | HTTP 框架，用于解析请求参数 |

**内部依赖**：
- `common/constants.go` — `ItemsPerPage` 变量（默认每页条数，值为 10）

## 3. 类型定义

| 类型名 | 字段 | 说明 |
|--------|------|------|
| `PageInfo` | `Page int` | 页码（从 1 开始） |
| | `PageSize int` | 每页条数 |
| | `Total int` | 总条数（后端设置） |
| | `Items any` | 返回数据（后端设置） |

## 4. 函数详解

### `(p *PageInfo) GetStartIdx() int`
返回分页查询的起始索引（0-based）。计算公式：`(Page - 1) * PageSize`。

### `(p *PageInfo) GetEndIdx() int`
返回分页查询的结束索引（不含）。计算公式：`Page * PageSize`。

### `(p *PageInfo) GetPageSize() int`
返回每页条数。

### `(p *PageInfo) GetPage() int`
返回当前页码。

### `(p *PageInfo) SetTotal(total int)`
设置总条数。

### `(p *PageInfo) SetItems(items any)`
设置返回数据。

### `GetPageQuery(c *gin.Context) *PageInfo`
从 Gin 请求上下文中解析分页参数并返回 `PageInfo` 指针。

**参数解析优先级**：
- 页码：优先读取 `p` 参数，若值 < 1 则重新解析 `p`（兼容 0 值），最终兜底为 1
- 页大小：优先读取 `page_size`，若为 0 则尝试 `ps`，再尝试 `size`（token page 场景），最终兜底为 `ItemsPerPage`

**边界处理**：
- 页码最小值为 1
- 页大小最大值为 100（超过则截断为 100）

## 5. 关键逻辑分析

1. **多参数名兼容**：为兼容不同前端/历史版本，支持 `page_size`、`ps`、`size` 三种参数名
2. **页码兜底**：当 `p` 解析失败或值 < 1 时，重新尝试解析，若仍为 0 则默认第 1 页
3. **页大小上限**：硬编码最大 100 条/页，防止客户端请求过大分页导致性能问题
4. **Total/Items 后设置**：`PageInfo` 作为响应体时，`Total` 和 `Items` 由业务逻辑层填充，不从请求中解析

## 6. 关联文件

- `common/constants.go` — 定义 `ItemsPerPage` 默认值（10）
- `controller/` — 各控制器调用 `GetPageQuery` 解析分页参数
- `model/` — 数据层使用 `GetStartIdx()`/`GetEndIdx()` 进行数据库分页查询
