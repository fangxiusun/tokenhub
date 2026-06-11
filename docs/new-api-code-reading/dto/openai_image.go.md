# openai_image.go 代码阅读文档

## 1. 全局摘要

该文件定义了 OpenAI 图像生成 API 的请求和响应数据结构。包含图像请求 `ImageRequest`（支持 DALL-E 等模型）、图像响应 `ImageResponse`，以及自定义的 JSON 序列化/反序列化逻辑，用于处理额外字段。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `reflect`：反射（用于字段名提取）
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数
  - `github.com/QuantumNous/new-api/types`：`TokenCountMeta` 类型
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### ImageRequest 结构体
图像生成请求结构：
- `Model` (string)：模型名称（如 dall-e-3）
- `Prompt` (string)：提示词（必填）
- `N` (*uint)：生成数量
- `Size` (string)：图像尺寸
- `Quality` (string)：质量等级
- `ResponseFormat` (string)：响应格式
- `Style` (json.RawMessage)：风格
- `User` (json.RawMessage)：用户标识
- `Background` (json.RawMessage)：背景
- `Moderation` (json.RawMessage)：审核设置
- `OutputFormat` (json.RawMessage)：输出格式
- `OutputCompression` (json.RawMessage)：输出压缩
- `PartialImages` (json.RawMessage)：部分图像
- `Stream` (*bool)：流式输出
- `Images` (json.RawMessage)：图像数组
- `Mask` (json.RawMessage)：遮罩
- `InputFidelity` (json.RawMessage)：输入保真度
- `Watermark` (*bool)：水印
- `Extra` (map[string]json.RawMessage)：额外字段（反序列化时捕获）

### ImageResponse 结构体
图像生成响应：
- `Data` ([]ImageData)：图像数据数组
- `Created` (int64)：创建时间戳
- `Metadata` (json.RawMessage)：元数据

### ImageData 结构体
图像数据：
- `Url` (string)：图像 URL
- `B64Json` (string)：Base64 编码图像数据
- `RevisedPrompt` (string)：修订后的提示词

## 4. 函数详情

### UnmarshalJSON()
```go
func (i *ImageRequest) UnmarshalJSON(data []byte) error
```
**功能**：自定义 JSON 反序列化，支持额外字段。

**逻辑**：
1. 解析原始 JSON 为 map
2. 提取已知字段名
3. 正常解析已知字段
4. 将未知字段存入 `Extra` map

### MarshalJSON()
```go
func (r ImageRequest) MarshalJSON() ([]byte, error)
```
**功能**：自定义 JSON 序列化。

**逻辑**：
1. 将结构体序列化为 map
2. 返回 map 的 JSON 表示

**注意**：注释掉的代码说明不合并 `ExtraFields`，避免重复。

### GetJSONFieldNames()
```go
func GetJSONFieldNames(t reflect.Type) map[string]struct{}
```
**功能**：通过反射获取结构体的 JSON 字段名。

**逻辑**：遍历结构体字段，提取 `json` 标签中的字段名。

### GetTokenCountMeta()
```go
func (i *ImageRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取 token 计数元数据。

**逻辑**：
1. 根据模型名称和尺寸计算尺寸比例
2. 根据模型和质量计算质量比例
3. 返回包含提示词和价格比例的元数据

**价格计算**：
- DALL-E 尺寸比例：256x256 (0.4), 512x512 (0.45), 1024x1024 (1.0), 1024x1792 (2.0)
- DALL-E-3 HD 质量比例：2.0（1024x1792 为 1.5）

### IsStream()
```go
func (i *ImageRequest) IsStream(c *gin.Context) bool
```
**功能**：判断是否为流式请求。

### SetModelName()
```go
func (i *ImageRequest) SetModelName(modelName string)
```
**功能**：设置模型名称。

## 5. 关键逻辑分析

1. **额外字段处理**：通过自定义序列化/反序列化，支持未定义字段的透传，保持与上游 API 的兼容性。

2. **价格计算逻辑**：根据模型、尺寸、质量动态计算图像价格比例，支持不同模型的定价策略。

3. **零值安全**：可选参数使用指针类型，确保序列化时能正确处理零值。

4. **字段名提取**：使用反射动态提取结构体字段名，避免硬编码字段列表。

## 6. 相关文件

- `relay/image/`：图像生成中继适配器
- `controller/image.go`：图像控制器
- `types/image.go`：图像相关类型定义
- `dto/usage.go`：`Usage` 结构体定义