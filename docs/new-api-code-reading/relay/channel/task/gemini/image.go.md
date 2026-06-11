# image.go 代码阅读文档

## 1. 全局总结

该文件提供 Gemini Veo 图片输入的解析工具函数，支持从 multipart 表单上传文件和 base64/data URI 字符串两种方式提取图片数据，并将其转换为 Veo API 所需的 `VeoImageInput` 格式。

## 2. 依赖关系

**标准库：**
- `encoding/base64` — Base64 编解码
- `io` — IO 操作
- `net/http` — MIME 类型检测
- `strings` — 字符串前缀匹配

**项目内部依赖：**
- `github.com/QuantumNous/new-api/constant` — 任务动作常量（TaskActionGenerate）
- `github.com/QuantumNous/new-api/relay/common` — RelayInfo 类型

**第三方依赖：**
- `github.com/gin-gonic/gin` — Web 框架上下文

## 3. 类型定义

### 常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `maxVeoImageSize` | 20 * 1024 * 1024 (20MB) | Veo 图片最大文件大小限制 |

## 4. 函数详解

| 函数签名 | 说明 |
|----------|------|
| `ExtractMultipartImage(c *gin.Context, info *relaycommon.RelayInfo) *VeoImageInput` | 从 multipart 表单中读取 `input_reference` 文件，返回 base64 编码的图片输入 |
| `ParseImageInput(imageStr string) *VeoImageInput` | 解析图片字符串（data URI 或纯 base64），返回 VeoImageInput |
| `parseDataURI(uri string) *VeoImageInput` | 解析 data URI 格式（如 `data:image/png;base64,...`） |

## 5. 关键逻辑分析

### multipart 文件提取流程
1. 从 `gin.Context` 获取 multipart 表单
2. 查找名为 `input_reference` 的文件字段
3. 检查文件大小是否超过 20MB 限制
4. 读取文件内容，检测 MIME 类型（优先使用文件头的 Content-Type，否则使用 `http.DetectContentType`）
5. 将文件内容 base64 编码后返回
6. 同时设置 `info.Action = TaskActionGenerate`（标记为有图片的生成任务）

### 图片字符串解析
- **data URI 格式**：解析 `data:{mime};base64,{data}` 格式，提取 MIME 类型和 base64 数据
- **纯 base64 格式**：直接解码验证有效性，使用 `http.DetectContentType` 检测 MIME 类型
- **空输入**：返回 nil

### 错误处理策略
所有解析函数在遇到错误时返回 nil 而非错误值，调用方需自行判断是否有图片输入。这是一种"尽力而为"的设计模式。

## 6. 关联文件

- `relay/channel/task/gemini/adaptor.go` — 调用 `ExtractMultipartImage` 和 `ParseImageInput`
- `relay/channel/task/gemini/dto.go` — `VeoImageInput` 结构体定义
