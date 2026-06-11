# notify-limit.go 代码阅读文档

## 1. 全局总结

该文件实现通知频率限制功能，支持 Redis 和内存两种存储方式。防止用户在短时间内收到过多通知。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | Redis 操作 |
| `constant` | 限制参数 |
| `gopool` | 后台清理任务 |

## 3. 类型定义

### `limitCount`
内存限制计数：`Count`、`Timestamp`

## 4. 函数详解

### `CheckNotificationLimit(userId, notifyType) (bool, error)`
统一入口，根据 Redis 是否启用选择存储方式

### `checkRedisLimit(userId, notifyType) (bool, error)`
Redis 实现：
- Key 格式：`notify_limit:{userId}:{notifyType}:{yyyyMMddHH}`
- 原子递增

### `checkMemoryLimit(userId, notifyType) (bool, error)`
内存实现：
- 使用 sync.Map 存储
- 启动后台清理任务（每小时）
- 自动过期清理

## 5. 关键逻辑分析

1. **时间窗口**：按小时计算（格式化到小时）
2. **自动清理**：内存模式下每小时清理过期条目
3. **原子操作**：Redis 使用 INCR 保证原子性

## 6. 关联文件

- `user_notify.go` — 通知发送
