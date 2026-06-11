# 类型层详细设计 (`types/`)

## 1. 概述
类型层包含整个应用程序中使用的共享类型定义。这些类型确保了跨包的类型一致性。

## 2. 文件详细说明

### 2.1 `error.go` - 错误类型
**职责**: 应用程序的自定义错误类型。

**关键类型**:
- **`NewAPIError`**: 自定义错误类型，带状态码和消息。
  - `Code int` - 错误代码
  - `Message string` - 错误消息
  - `Error() string` - 实现 error 接口
- **`NewAPIErrorWithStatusCode`**: 带 HTTP 状态码的错误。
  - `StatusCode int` - HTTP 状态码
  - `Error *NewAPIError` - 底层错误

### 2.2 `relay_format.go` - 中转格式类型
**职责**: 中转格式类型和常量。

**关键类型**:
- **`RelayFormat`**: 中转格式类型。
- **`RelayFormatOpenAI`**: OpenAI 格式常量。
- **`RelayFormatClaude`**: Claude 格式常量。
- **`RelayFormatGemini`**: Gemini 格式常量。

### 2.3 `price_data.go` - 价格数据
**职责**: 计费的价格数据结构。

**关键类型**:
- **`PriceData`**: 包含模型的定价信息。
  - `ModelRatio float64` - 模型比例
  - `CompletionRatio float64` - 完成比例
  - `CacheRatio float64` - 缓存比例
  - `AudioRatio float64` - 音频比例

### 2.4 `file_source.go` - 文件源类型
**职责**: 媒体处理的文件源类型。

**关键类型**:
- **`FileSource`**: 表示文件源（URL, base64 等）。
  - `Type string` - 文件源类型
  - `Data string` - 文件数据
  - `URL string` - 文件 URL

## 3. 使用场景
- **错误处理**: 在控制器和服务层使用 `NewAPIError` 进行错误处理
- **中转格式**: 在中转层使用 `RelayFormat` 跟踪格式转换
- **计费**: 在计费系统中使用 `PriceData` 进行价格计算
- **文件处理**: 在图像/音频处理中使用 `FileSource` 处理文件源

---

## 关联文件列表

### 类型层核心文件
- `types/error.go` - 错误类型
- `types/relay_format.go` - 中转格式类型
- `types/price_data.go` - 价格数据
- `types/file_source.go` - 文件源类型

### 依赖此模块的文件
- `controller/relay.go` - 中转控制器（使用 NewAPIError）
- `relay/channel/adapter.go` - 适配器接口（使用 RelayFormat）
- `service/billing.go` - 计费逻辑（使用 PriceData）
- `dto/error.go` - 错误 DTO（使用 NewAPIError）
- `dto/openai_request.go` - OpenAI 请求（使用 FileSource）

### 依赖的外部库
- 无外部依赖，纯类型定义
