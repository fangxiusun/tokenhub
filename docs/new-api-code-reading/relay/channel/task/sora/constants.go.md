# constants.go 代码阅读文档

## 1. 全局总结

本文件定义了 Sora 视频生成通道的常量配置，包括支持的模型列表和通道名称标识。这些常量被 Sora 适配器（`adaptor.go`）和其他引用 Sora 通道的代码所使用，是 Sora 通道注册和路由的基础元数据。

## 2. 依赖关系

无外部依赖，纯定义文件。

## 3. 类型定义

### 变量

| 名称 | 类型 | 值 | 说明 |
|------|------|-----|------|
| `ModelList` | `[]string` | `["sora-2", "sora-2-pro"]` | Sora 通道支持的模型列表，用于模型白名单校验和前端模型选择展示 |
| `ChannelName` | `string` | `"sora"` | 通道名称标识，用于日志、错误信息和通道注册 |

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

- `ModelList` 定义了两个模型：`sora-2`（基础版）和 `sora-2-pro`（专业版），对应 OpenAI Sora 视频生成服务的不同级别。
- `ChannelName` 作为通道的唯一字符串标识，在适配器的 `GetChannelName()` 方法中返回。
- 该文件是典型的 Go 常量定义模式，将易变的配置集中管理，便于后续扩展新模型。

## 6. 关联文件

- `relay/channel/task/sora/adaptor.go` — Sora 适配器实现，引用 `ModelList` 和 `ChannelName`
- `relay/common/relay_info.go` — 通道信息结构体，使用通道类型常量
- `constant/` — 项目常量定义包，定义了通道类型常量
