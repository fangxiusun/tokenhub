# dto.go 代码阅读文档

## 1. 全局总结
本文件定义了 Replicate API 专用的数据传输对象，包括预测响应、预测错误和文件上传响应的结构体。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `PredictionResponse`
Replicate 预测响应结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Status` | `string` | `json:"status"` | 预测状态（如 "succeeded"、"processing"、"failed"） |
| `Output` | `any` | `json:"output"` | 预测输出（可以是字符串、数组或其他类型） |
| `Error` | `*PredictionError` | `json:"error"` | 错误信息（指针类型，无错误时为 nil） |

### `PredictionError`
Replicate 预测错误结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Code` | `string` | `json:"code"` | 错误码 |
| `Message` | `string` | `json:"message"` | 错误消息 |
| `Detail` | `string` | `json:"detail"` | 详细错误信息 |

### `FileUploadResponse`
Replicate 文件上传响应结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Urls` | `struct{ Get string }` | `json:"urls"` | 文件 URL 信息 |
| `Urls.Get` | `string` | `json:"get"` | 文件的 GET URL |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- **Output 类型灵活**：`PredictionResponse.Output` 使用 `any` 类型，因为 FLUX 模型的输出可以是单个 URL 字符串或 URL 数组。
- **Error 指针类型**：`Error` 使用指针类型，允许区分"无错误"（nil）和"空错误对象"。
- **FileUploadResponse 嵌套结构**：使用匿名结构体嵌套 `Urls.Get`，匹配 Replicate API 的 JSON 结构。

## 6. 关联文件
- `relay/channel/replicate/adaptor.go` — 使用这些结构体进行响应解析
- `relay/channel/replicate/constants.go` — 渠道常量
