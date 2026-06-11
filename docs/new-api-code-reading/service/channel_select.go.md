# channel_select.go 代码阅读文档

## 1. 全局总结

该文件实现了通道选择的核心逻辑，特别是自动分组（Auto Group）的通道选择策略。支持跨分组重试、优先级遍历、以及分组切换的复杂流程控制。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 重试次数、上下文操作 |
| `constant` | 上下文键名常量 |
| `logger` | 调试日志 |
| `model` | 通道查询 |
| `setting` | 自动分组配置 |
| `gin` | HTTP 上下文 |

## 3. 类型定义

### `RetryParam`
重试参数结构体：
- `Ctx` — HTTP 上下文
- `TokenGroup` — 令牌分组
- `ModelName` — 模型名称
- `Retry` — 当前重试次数（指针，支持外部修改）
- `resetNextTry` — 下次重试时重置计数器标志

## 4. 函数详解

### `CacheGetRandomSatisfiedChannel(param) (*model.Channel, string, error)`
核心通道选择函数，支持两种模式：

**普通模式**：直接调用 `model.GetRandomSatisfiedChannel`

**自动分组模式**（`param.TokenGroup == "auto"`）：
1. 获取用户可用的自动分组列表
2. 从上次分组索引开始遍历
3. 每个分组内按优先级重试
4. 当前分组优先级用完时，切换到下一个分组
5. 跨分组重试时重置重试计数器

### `RetryParam` 方法
- `GetRetry()` — 获取当前重试次数
- `SetRetry(retry)` — 设置重试次数
- `IncreaseRetry()` — 增加重试次数（支持 resetNextTry）
- `ResetRetryNextTry()` — 标记下次重试时重置

## 5. 关键逻辑分析

1. **分组切换机制**：通过 `ContextKeyAutoGroupIndex` 跟踪当前分组
2. **优先级重试**：每个分组内独立的优先级遍历
3. **跨分组重试**：当前分组用完后自动切换到下一个分组
4. **状态持久化**：通过 gin.Context 在重试间传递状态

## 6. 关联文件

- `channel_affinity.go` — 通道亲和性
- `model/channel.go` — GetRandomSatisfiedChannel
- `setting` — 自动分组配置
