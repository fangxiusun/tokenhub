# ability.go 代码阅读文档

## 1. 全局总结
ability.go 定义了渠道能力模型，存储 (group, model, channel_id) 映射关系，是渠道路由的核心数据结构。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `service/channel.go` - 渠道缓存管理
- `middleware/distributor.go` - 渠道路由分发
- `controller/channel.go` - 渠道管理

## 3. 类型定义
### 3.1 结构体
- `Ability` - 渠道能力结构体
  - `Group string` - 用户组
  - `Model string` - 模型名称
  - `ChannelId int` - 渠道 ID
  - `Enabled bool` - 是否启用
  - `Priority *int64` - 优先级
  - `Weight uint` - 权重

## 4. 函数详解
### 4.1 CreateAbility
- **职责**: 创建能力记录

### 4.2 DeleteAbilitiesByChannelId
- **职责**: 删除指定渠道的所有能力

### 4.3 GetRandomSatisfiedChannel
- **职责**: 根据模型和组随机选择满足条件的渠道

## 5. 关键逻辑分析
- **优先级选择**: 按优先级降序选择，同优先级内按权重随机
- **批量更新**: 支持批量创建和删除能力

## 6. 关联文件
- `service/channel_select.go` - 渠道选择算法
- `middleware/distributor.go` - 渠道路由分发
