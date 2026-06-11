# multi_key_mode.go 代码阅读文档

## 1. 全局概述

本文件定义了多密钥模式（MultiKeyMode）常量，用于控制当一个渠道配置了多个 API Key 时的密钥选择策略。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### MultiKeyMode 类型

```go
type MultiKeyMode string
```

字符串类型的多密钥模式标识符。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 多密钥模式常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `MultiKeyModeRandom` | `"random"` | 随机选择密钥 |
| `MultiKeyModePolling` | `"polling"` | 轮询选择密钥 |

### 使用场景

当一个渠道配置了多个 API Key 时（如 OpenAI 渠道配置了多个 Key 以提高并发能力），系统需要决定每次请求使用哪个 Key：

- **随机模式**：每次请求随机选择一个 Key，适用于负载均衡场景
- **轮询模式**：按顺序依次使用 Key，适用于需要均匀分配请求的场景

### 实现位置

密钥选择逻辑在 `model/channel.go` 和 `relay/` 相关代码中实现，本文件仅提供常量定义。

## 6. 相关文件

- `model/channel.go` — 渠道模型，包含多密钥存储和选择逻辑
- `relay/relay_info.go` — 中继信息中使用多密钥模式
- `middleware/distributor.go` — 渠道分发时使用多密钥模式
