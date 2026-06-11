# channel_error.go 代码阅读文档

## 1. 全局概述

本文件定义了渠道错误（ChannelError）结构体，用于在请求处理过程中记录出错渠道的详细信息，便于日志记录和问题排查。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### ChannelError 结构体

```go
type ChannelError struct {
    ChannelId   int    `json:"channel_id"`
    ChannelType int    `json:"channel_type"`
    ChannelName string `json:"channel_name"`
    IsMultiKey  bool   `json:"is_multi_key"`
    AutoBan     bool   `json:"auto_ban"`
    UsingKey    string `json:"using_key"`
}
```

| 字段 | 说明 |
|------|------|
| `ChannelId` | 渠道 ID |
| `ChannelType` | 渠道类型 |
| `ChannelName` | 渠道名称 |
| `IsMultiKey` | 是否为多密钥渠道 |
| `AutoBan` | 是否启用自动封禁 |
| `UsingKey` | 当前使用的密钥标识 |

## 4. 函数详情

### NewChannelError

```go
func NewChannelError(channelId int, channelType int, channelName string, isMultiKey bool, usingKey string, autoBan bool) *ChannelError
```

创建新的 ChannelError 实例。

## 5. 关键逻辑分析

### 设计目的

`ChannelError` 的主要作用是在请求失败时记录出错渠道的上下文信息：
- 便于定位是哪个渠道出了问题
- 记录使用的密钥标识，便于排查密钥问题
- `AutoBan` 信息用于决定是否自动禁用该渠道

### 使用场景

- 渠道请求失败时，将 `ChannelError` 附加到 `NewAPIError` 中
- 错误日志中记录渠道信息
- 自动封禁逻辑根据 `AutoBan` 决定是否禁用渠道

## 6. 相关文件

- `types/error.go` — NewAPIError 中可能包含 ChannelError
- `model/channel.go` — 渠道模型，使用 ChannelError 进行自动封禁
- `relay/` — 中继层在请求失败时创建 ChannelError
