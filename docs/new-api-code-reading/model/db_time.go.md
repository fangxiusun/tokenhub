# db_time.go 代码阅读文档

## 1. 全局总结
db_time.go 提供数据库时间处理工具函数，处理不同数据库的时间格式差异。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具

### 2.2 被引用的文件
- `model/` - 模型层时间处理

## 3. 类型定义
### 3.1 函数
- `GetDBTime() int64` - 获取数据库当前时间
- `GetDBTimestamp() string` - 获取数据库时间戳字符串

## 4. 函数详解
### 4.1 GetDBTime
- **签名**: `func GetDBTime() int64`
- **职责**: 获取数据库服务器当前时间（Unix 时间戳）

## 5. 关键逻辑分析
- **跨数据库兼容**: 处理 SQLite、MySQL、PostgreSQL 的时间格式差异

## 6. 关联文件
- `common/database.go` - 数据库工具
