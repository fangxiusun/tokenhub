# main.go (model) 代码阅读文档

## 1. 全局总结
main.go 是模型层的初始化文件，负责数据库连接建立、表结构迁移、批量更新支持等核心功能。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具（数据库、Redis）
- `constant` - 常量定义
- `gorm.io/gorm` - ORM
- `gorm.io/driver/sqlite` - SQLite 驱动
- `gorm.io/driver/mysql` - MySQL 驱动
- `gorm.io/driver/postgres` - PostgreSQL 驱动

### 2.2 被引用的文件
- `main.go` - 应用入口

## 3. 类型定义
### 3.1 变量
- `DB` - 全局数据库实例
- `batchUpdates` - 批量更新通道

## 4. 函数详解
### 4.1 Init
- **职责**: 初始化数据库连接
- **逻辑流程**:
  1. 根据环境变量选择数据库类型
  2. 建立数据库连接
  3. 执行表结构迁移
  4. 初始化批量更新

### 4.2 GetDB
- **职责**: 获取数据库实例

### 4.3 Migration
- **职责**: 执行数据库迁移

## 5. 关键逻辑分析
- **多数据库支持**: 支持 SQLite、MySQL、PostgreSQL
- **批量更新**: 使用通道实现异步批量写入
- **自动迁移**: 应用启动时自动执行表结构迁移

## 6. 关联文件
- `common/database.go` - 数据库工具
- `main.go` - 应用入口
