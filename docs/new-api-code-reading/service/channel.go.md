# channel.go 代码阅读文档

## 1. 全局总结

该文件提供通道（Channel）管理的核心业务逻辑，包括通道的自动禁用、启用、以及判断是否应该禁用/启用通道的决策函数。这些函数被 relay 层在请求失败时调用，实现通道的自动健康管理。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 通道状态常量、自动禁用开关 |
| `dto` | 通知类型常量 |
| `model` | 通道状态更新 |
| `operation_setting` | 自动禁用配置 |
| `types` | 错误类型判断 |

## 3. 类型定义

无自定义类型，仅使用已有类型。

## 4. 函数详解

### `formatNotifyType(channelId, status) string`
生成通知类型标识字符串，格式：`channel_update_{channelId}_{status}`

### `DisableChannel(channelError, reason)`
- 检查是否启用自动禁用（`AutoBan` 标志）
- 调用 `model.UpdateChannelStatus` 更新为自动禁用状态
- 成功后通知管理员

### `EnableChannel(channelId, usingKey, channelName)`
- 启用通道并通知管理员

### `ShouldDisableChannel(err) bool`
判断是否应该禁用通道的决策函数：
1. 自动禁用功能未启用 → false
2. 是 ChannelError → true
3. 是 SkipRetryError → false
4. HTTP 状态码匹配禁用规则 → true
5. 错误消息包含自动禁用关键词 → true

### `ShouldEnableChannel(newAPIError, status) bool`
判断是否应该启用通道：
1. 自动启用功能未启用 → false
2. 有错误 → false
3. 状态不是自动禁用 → false
4. 以上条件都满足 → true

## 5. 关键逻辑分析

1. **双层判断**：ShouldDisableChannel 先检查错误类型，再检查状态码，最后检查关键词
2. **AC自动机搜索**：使用 Aho-Corasick 算法进行关键词匹配
3. **通知去重**：通过 formatNotifyType 生成唯一通知标识

## 6. 关联文件

- `str.go` — AcSearch 自动机搜索实现
- `model/channel.go` — 通道数据库操作
- `setting/operation_setting` — 自动禁用配置
