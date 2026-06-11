# log.go 代码阅读文档

## 1. 全局总结
log.go 定义了日志模型，记录 API 请求的详细信息，包括用户、渠道、Token 使用量等。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `controller/log.go` - 日志控制器
- `service/` - 日志记录

## 3. 类型定义
### 3.1 结构体
- `Log` - 日志结构体
  - `Id int` - 日志 ID
  - `UserId int` - 用户 ID
  - `CreatedAt int64` - 创建时间
  - `Type int` - 日志类型
  - `Content string` - 日志内容
  - `ModelName string` - 模型名称
  - `TokenName string` - Token 名称
  - `Quota int` - 消耗配额
  - `ChannelId int` - 渠道 ID
  - `Channel string` - 渠道名称
  - `PromptTokens int` - 输入 Token 数
  - `CompletionTokens int` - 输出 Token 数

## 4. 函数详解
### 4.1 CreateLog
- **职责**: 创建日志记录

### 4.2 GetLogsByUserId
- **职责**: 获取用户日志

### 4.3 GetAllLogs
- **职责**: 获取所有日志

### 4.4 GetLogsStat
- **职责**: 获取日志统计

## 5. 关键逻辑分析
- **配额跟踪**: 记录每次请求消耗的配额
- **渠道跟踪**: 记录使用的渠道信息
- **Token 统计**: 记录输入输出 Token 数

## 6. 关联文件
- `controller/log.go` - 日志控制器
- `service/` - 日志记录服务
