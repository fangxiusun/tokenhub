# database.go 代码阅读文档

## 1. 全局总结

该文件定义了项目中数据库类型相关的常量和全局变量，用于在运行时标识和切换不同的数据库后端（SQLite、MySQL、PostgreSQL、ClickHouse）。文件提供了简洁的数据库类型判断机制，供其他模块根据当前使用的数据库类型执行相应的逻辑。

## 2. 依赖关系

本文件无外部依赖，仅使用 Go 内置功能。

## 3. 类型定义

### 3.1 数据库类型常量

```go
const (
    DatabaseTypeMySQL      = "mysql"
    DatabaseTypeSQLite     = "sqlite"
    DatabaseTypePostgreSQL = "postgres"
)
```

**用途**：定义三种主要数据库类型的字符串标识，用于日志、配置和条件判断。

**注意事项**：
- `DatabaseTypePostgreSQL` 的值为 `"postgres"` 而非 `"postgresql"`，这与 GORM 的命名约定一致
- ClickHouse 未定义为常量，但有对应的全局变量 `UsingClickHouse`

## 4. 变量详解

### 4.1 数据库类型标志变量

| 变量名 | 类型 | 默认值 | 用途 |
|--------|------|--------|------|
| `UsingSQLite` | `bool` | `false` | 标识是否使用 SQLite |
| `UsingPostgreSQL` | `bool` | `false` | 标识是否使用 PostgreSQL |
| `UsingMySQL` | `bool` | `false` | 标识是否使用 MySQL |
| `UsingClickHouse` | `bool` | `false` | 标识是否使用 ClickHouse |

**使用场景**：
- 在数据库迁移时判断是否需要执行特定的 SQL 语法
- 在查询构建时选择正确的列引用方式
- 在日志记录时标识数据库类型

### 4.2 LogSqlType 变量

```go
var LogSqlType = DatabaseTypeSQLite
```

**用途**：指定 SQL 日志记录时使用的数据库类型，默认为 SQLite。

**设计考虑**：
- 该变量独立于实际使用的数据库类型
- 允许在不改变实际数据库的情况下调整日志格式

### 4.3 SQLitePath 变量

```go
var SQLitePath = "one-api.db?_busy_timeout=30000"
```

**用途**：定义 SQLite 数据库文件路径和连接参数。

**参数说明**：
- `one-api.db`: 数据库文件名
- `_busy_timeout=30000`: 设置忙等待超时为 30 秒，避免并发写入时的锁冲突

## 5. 函数详解

本文件无函数定义，仅包含常量和变量声明。

## 6. 关键逻辑分析

### 6.1 多数据库支持机制

项目采用运行时标志变量（`UsingSQLite`、`UsingMySQL` 等）来支持多数据库：
- 在应用启动时根据配置设置对应的标志
- 在业务代码中通过 `if common.UsingMySQL` 等条件分支处理数据库差异

### 6.2 数据库类型常量的用途

这些常量主要用于：
- **日志输出**：在 SQL 日志中标识数据库类型
- **配置管理**：在配置文件中使用统一的字符串标识
- **GORM 兼容**：与 GORM 框架的数据库类型标识保持一致

### 6.3 ClickHouse 的特殊处理

ClickHouse 作为分析型数据库，与事务型数据库（SQLite、MySQL、PostgreSQL）有本质区别：
- 仅用于日志分析和统计查询
- 不参与主业务逻辑
- 因此没有对应的数据库类型常量

### 6.4 SQLite 忙等待配置

`_busy_timeout=30000` 参数解决了 SQLite 在并发场景下的常见问题：
- SQLite 使用文件锁，写入时会阻塞其他写入操作
- 30 秒的超时时间允许等待其他操作完成，而不是立即返回错误
- 适用于中等并发场景

## 7. 关联文件

- `model/main.go`：数据库初始化，设置这些全局变量
- `setting/` 目录：数据库配置管理
- `controller/` 目录：根据数据库类型执行不同的查询逻辑
- `common/mysql.go`、`common/postgres.go`、`common/sqlite.go`：各数据库的特定工具函数
