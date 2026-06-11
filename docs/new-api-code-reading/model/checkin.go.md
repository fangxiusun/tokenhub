# checkin.go 代码阅读文档

## 1. 全局总结
checkin.go 实现了每日签到功能，记录用户签到状态和连续签到天数。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `gorm.io/gorm` - ORM

### 2.2 被引用的文件
- `controller/checkin.go` - 签到控制器
- `service/` - 签到奖励逻辑

## 3. 类型定义
### 3.1 结构体
- `CheckIn` - 签到记录结构体
  - `Id int` - 记录 ID
  - `UserId int` - 用户 ID
  - `ConsecutiveDays int` - 连续签到天数
  - `LastCheckInTime int64` - 最后签到时间

## 4. 函数详解
### 4.1 GetCheckInStatus
- **职责**: 获取用户签到状态

### 4.2 DoCheckIn
- **职责**: 执行签到操作

## 5. 关键逻辑分析
- **连续签到**: 自动计算连续签到天数
- **奖励计算**: 根据连续天数计算奖励

## 6. 关联文件
- `controller/checkin.go` - 签到控制器
