# channel.go 代码阅读文档

## 1. 全局总结
channel.go 定义了渠道模型，是系统的核心数据结构之一。包含渠道的完整信息：类型、密钥、基础 URL、模型列表、状态等。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `service/channel.go` - 渠道缓存管理
- `controller/channel.go` - 渠道管理
- `relay/relay_adaptor.go` - 适配器工厂

## 3. 类型定义
### 3.1 结构体
- `Channel` - 渠道结构体
  - `Id int` - 渠道 ID
  - `Type int` - 渠道类型
  - `Key string` - API 密钥
  - `BaseURL string` - 基础 URL
  - `Status int` - 状态
  - `Name string` - 渠道名称
  - `Weight int` - 权重
  - `Models string` - 支持的模型列表
  - `ModelMapping string` - 模型映射
  - `groups string` - 用户组
  - `tags string` - 标签
  - `AutoBan bool` - 自动禁用
  - `Balance float64` - 余额
  - `BalanceUpdatedTime int64` - 余额更新时间

## 4. 函数详解
### 4.1 GetAllChannels
- **职责**: 获取所有渠道

### 4.2 CreateChannel
- **职责**: 创建渠道

### 4.3 UpdateChannel
- **职责**: 更新渠道

### 4.4 DeleteChannel
- **职责**: 删除渠道

### 4.5 GetAvailableChannels
- **职责**: 获取可用渠道列表

## 5. 关键逻辑分析
- **多密钥支持**: Key 字段支持逗号分隔的多个密钥
- **模型映射**: ModelMapping 支持模型名称别名
- **状态管理**: 支持启用/禁用/自动禁用

## 6. 关联文件
- `model/ability.go` - 渠道能力
- `service/channel.go` - 渠道缓存
- `relay/relay_adaptor.go` - 适配器工厂
