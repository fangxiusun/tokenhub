# helpers.go 代码阅读文档

## 1. 全局总结

本文件是任务子系统的通用工具集，提供了多个适配器共享的辅助函数和基础结构体。核心功能包括：元数据反序列化、默认值处理、任务 ID 编解码、代理 URL 构建、进度常量定义，以及可嵌入的无操作计费基础结构体 `BaseBilling`。这些工具被所有任务适配器（Sora、Suno、Vertex、Vidu 等）广泛使用。

## 2. 依赖关系

### 标准库
- `encoding/base64` — Base64 编解码
- `fmt` — 格式化输出

### 项目内部包
- `github.com/QuantumNous/new-api/common` — JSON 操作
- `github.com/QuantumNous/new-api/model` — 任务模型类型
- `relaycommon`（别名）— 中继公共结构体
- `github.com/QuantumNous/new-api/setting/system_setting` — 系统设置（服务器地址）

### 第三方库
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### 结构体

| 名称 | 字段 | 说明 |
|------|------|------|
| `BaseBilling` | 无字段 | 空结构体，提供 `EstimateBilling`、`AdjustBillingOnSubmit`、`AdjustBillingOnComplete` 的默认无操作实现 |

### 常量

| 名称 | 值 | 说明 |
|------|-----|------|
| `ProgressSubmitted` | `"10%"` | 已提交进度 |
| `ProgressQueued` | `"20%"` | 排队中进度 |
| `ProgressInProgress` | `"30%"` | 进行中进度 |
| `ProgressComplete` | `"100%"` | 完成进度 |

## 4. 函数详解

### `UnmarshalMetadata(metadata map[string]any, target any) error`
- **作用**: 将 `map[string]any` 类型的元数据通过 JSON 序列化/反序列化转换为指定类型结构体。
- **安全特性**: 会删除 metadata 中的 `model` 字段，防止元数据覆盖模型字段导致计费绕过。
- **参数**: `metadata` — 原始元数据 map；`target` — 目标结构体指针。
- **返回**: 序列化或反序列化失败时返回错误。

### `DefaultString(val, fallback string) string`
- **作用**: 如果 `val` 为空字符串，返回 `fallback`，否则返回 `val`。
- **用途**: 为请求参数提供默认值。

### `DefaultInt(val, fallback int) int`
- **作用**: 如果 `val` 为 0，返回 `fallback`，否则返回 `val`。
- **用途**: 为整数请求参数提供默认值。

### `EncodeLocalTaskID(name string) string`
- **作用**: 将上游操作名称（如 Vertex 的 `projects/.../operations/...`）编码为 URL 安全的 Base64 字符串。
- **编码方式**: `base64.RawURLEncoding`（无填充、URL 安全字符集）。
- **用途**: Vertex/Gemini 适配器将上游操作名编码为本地任务 ID。

### `DecodeLocalTaskID(id string) (string, error)`
- **作用**: 解码 Base64 编码的任务 ID，还原为上游操作名称。
- **返回**: 解码后的字符串和可能的错误。

### `BuildProxyURL(taskID string) string`
- **作用**: 构建视频代理 URL，格式为 `{ServerAddress}/v1/videos/{taskID}/content`。
- **用途**: 为客户端提供统一的视频内容访问入口。

## 5. 关键逻辑分析

### BaseBilling 设计模式
`BaseBilling` 采用 Go 的结构体嵌入模式，提供"零配置"的默认实现：
- `EstimateBilling` → 返回 `nil`（使用基础模型价格）
- `AdjustBillingOnSubmit` → 返回 `nil`（提交时无调整）
- `AdjustBillingOnComplete` → 返回 `0`（保持预扣金额）

适配器只需嵌入 `BaseBilling` 即可满足接口要求，需要自定义计费时重写对应方法（如 Sora 重写了 `EstimateBilling`）。

### 元数据安全
`UnmarshalMetadata` 中删除 `model` 字段是重要的安全措施。如果攻击者在元数据中注入 `model` 字段，可能绕过基于模型的计费规则。

### 任务 ID 编码
`EncodeLocalTaskID`/`DecodeLocalTaskID` 使用 Base64 Raw URL 编码，适用于 URL 路径和 JSON 字段，避免了 `+`、`/`、`=` 等特殊字符的问题。

## 6. 关联文件

- `relay/channel/task/sora/adaptor.go` — 使用 `BaseBilling`、`DefaultString`、`DefaultInt`
- `relay/channel/task/suno/adaptor.go` — 使用 `BaseBilling`
- `relay/channel/task/vertex/adaptor.go` — 使用 `BaseBilling`、`EncodeLocalTaskID`、`DecodeLocalTaskID`
- `relay/channel/task/vidu/adaptor.go` — 使用 `BaseBilling`、`DefaultString`、`DefaultInt`、`UnmarshalMetadata`
- `relay/common/relay_info.go` — `TaskSubmitReq` 结构体定义
- `setting/system_setting/` — `ServerAddress` 系统设置
