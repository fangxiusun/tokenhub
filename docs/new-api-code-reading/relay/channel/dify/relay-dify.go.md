# relay-dify.go 代码阅读文档

## 1. 全局总结

本文件实现了 Dify 频道的请求转换和响应处理逻辑。包含文件上传、OpenAI → Dify 请求转换、流式/非流式响应处理，以及 Dify 事件到 OpenAI SSE 格式的转换。支持图片文件的 base64 解码上传和远程 URL 引用。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `bytes` | 内存缓冲区 |
| `encoding/base64` | Base64 解码 |
| `encoding/json` | JSON 操作 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `mime/multipart` | multipart 表单构建 |
| `net/http` | HTTP 请求/响应 |
| `os` | 临时文件操作 |
| `strings` | 字符串处理 |
| `common` | JSON 工具、日志、时间戳 |
| `constant` | 全局常量（DifyDebug） |
| `dto` | 数据传输对象 |
| `relay/common` | RelayInfo 上下文 |
| `relay/helper` | 流式处理工具 |
| `service` | HTTP 客户端、Token 估算 |
| `types` | 错误类型 |
| `lo` | 工具库（FromPtrOr） |
| `gin` | Web 框架 |

## 3. 类型定义

无类型定义。

## 4. 函数详解

### `uploadDifyFile(c *gin.Context, info *relaycommon.RelayInfo, user string, media dto.MediaContent) *DifyFile`
上传文件到 Dify：
1. 仅处理 `ContentTypeImageURL` 类型
2. 从 base64 数据中解码图片（去除 `data:image/...;base64,` 前缀）
3. 创建临时文件写入解码数据
4. 构建 multipart 表单（包含 user 字段和 file 字段）
5. 发送 POST 请求到 `/v1/files/upload`
6. 解析响应获取文件 ID
7. 返回 `DifyFile`（type: "image", transfer_mode: "local_file"）

### `requestOpenAI2Dify(c *gin.Context, info *relaycommon.RelayInfo, request dto.GeneralOpenAIRequest) *DifyChatRequest`
将 OpenAI 请求转换为 Dify 格式：
1. 初始化空的 `Inputs` 和 `Files`
2. 处理 `User` 字段（默认使用 response ID）
3. 遍历消息列表：
   - `system` → `"SYSTEM: \n{content}\n"`
   - `assistant` → `"ASSISTANT: \n{content}\n"`
   - `user` → 解析内容，处理文本和图片
4. 图片处理：
   - 远程图片 → 创建 `DifyFile`（transfer_mode: "remote_url"）
   - 本地图片 → 调用 `uploadDifyFile` 上传
5. 拼接所有文本为 `Query`
6. 根据 stream 标志设置 `ResponseMode`（"blocking" / "streaming"）

### `streamResponseDify2OpenAI(difyResponse DifyChunkChatCompletionResponse) *dto.ChatCompletionsStreamResponse`
将 Dify 流式事件转换为 OpenAI SSE 格式：
- **workflow_\* 事件**: 调试模式下输出工作流 ID 和状态到 reasoning_content
- **node_\* 事件**: 调试模式下输出节点类型和状态到 reasoning_content
- **message / agent_message 事件**: 将 answer 作为 content 输出
  - 特殊处理 thinking 标签：将 HTML details 标签转换为 `<think>` 标签

### `difyStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Dify 流式响应：
1. 使用 `helper.StreamScannerHandler` 进行流式扫描
2. 对每个事件：
   - 解析为 `DifyChunkChatCompletionResponse`
   - `message_end` → 提取 usage，结束
   - `error` → 停止处理
   - 其他 → 转换为 OpenAI 格式并发送
3. 累加响应文本和节点 token 数
4. 流结束后，如果 usage 为空则估算
5. 将节点 token 数加到 completion_tokens

### `difyHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Dify 非流式响应：
1. 读取并解析响应体
2. 构建 `dto.OpenAITextResponse`
3. 设置 conversation_id 作为响应 ID
4. 序列化并返回

## 5. 关键逻辑分析

### 文件上传流程
```
base64 图片 → 解码 → 临时文件 → multipart 表单 → POST /v1/files/upload → 获取 file_id
```
这是一个同步阻塞操作，对于大量图片可能成为性能瓶颈。

### Thinking 标签转换
Dify 的思考链输出使用 HTML `<details>` 标签包裹，需要转换为标准的 `<think>` 标签：
```
<details open><summary>Thinking...</summary>\n → <think>
</details> → </think>
```

### 调试模式（DifyDebug）
当 `constant.DifyDebug` 为 true 时，工作流和节点事件的信息会输出到 `reasoning_content` 字段，方便调试复杂的 Dify 工作流。

### Token 用量的双重来源
1. 优先使用 Dify API 返回的 `metadata.usage`（精确值）
2. 回退到文本估算（近似值）
3. 额外计算节点 token 数（调试信息）

### 消息拼接策略
OpenAI 的多条消息被拼接为单一的 `Query` 文本，使用 `SYSTEM:`、`ASSISTANT:`、`USER:` 前缀区分角色。这丢失了原始的消息结构，但对于 Dify 的单轮对话模式是可接受的。

### 远程图片的处理
远程图片（http/https URL）直接作为 `DifyFile` 传递，不需要上传，减少了网络开销。

## 6. 关联文件

- `relay/channel/dify/adaptor.go` - 调度本文件中的函数
- `relay/channel/dify/dto.go` - DTO 定义
- `relay/helper/` - StreamScannerHandler、ObjectData 等工具
- `constant/` - DifyDebug 全局开关
- `service/` - HTTP 客户端、Token 估算
