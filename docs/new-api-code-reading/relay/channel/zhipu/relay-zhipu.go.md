# relay-zhipu.go 代码阅读文档

## 1. 全局总结
智谱 AI 旧版渠道的核心中继逻辑文件。包含 JWT Token 生成与缓存、请求/响应格式转换、流式/非流式处理器。

## 2. 依赖关系
- **标准库**: bufio, encoding/json, io, net/http, strings, sync, time
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理、日志、时间戳
  - `github.com/QuantumNous/new-api/constant` — 常量定义
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/helper` — 流式处理工具
  - `github.com/QuantumNous/new-api/service` — 响应体处理
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/golang-jwt/jwt/v5` — JWT 签发
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### 全局变量
```go
var zhipuTokens sync.Map  // JWT Token 缓存
var expSeconds int64 = 24 * 3600  // Token 有效期 24 小时
```

## 4. 函数详解

### getZhipuToken(apikey string) string
生成智谱 JWT Token：
1. 从缓存获取有效 Token
2. 解析 API Key（格式: `id.secret`）
3. 构建 JWT Claims（api_key、exp、timestamp）
4. 使用 HS256 算法签名
5. 设置 `sign_type: SIGN` 头
6. 缓存 Token 供后续使用

### requestOpenAI2Zhipu(request dto.GeneralOpenAIRequest) *ZhipuRequest
将 OpenAI 请求转换为智谱格式：
- 使用 `prompt` 字段（而非 `messages`）
- system 消息转换为 user+assistant 对话对
- 设置 Temperature、TopP、Incremental

### responseZhipu2OpenAI(response *ZhipuResponse) *dto.OpenAITextResponse
将智谱响应转换为 OpenAI 格式：
- 遍历 Choices 数组
- 最后一个 Choice 的 FinishReason 设为 "stop"
- Content 去除首尾引号

### streamResponseZhipu2OpenAI(zhipuResponse string) *dto.ChatCompletionsStreamResponse
将智谱流式数据转换为 OpenAI 流式格式。

### streamMetaResponseZhipu2OpenAI(...)
处理智谱流式响应的元数据（meta: 前缀），返回 Usage 信息。

### zhipuStreamHandler(...)
流式响应处理器：
1. 使用 goroutine 读取 SSE 流
2. 区分 `data:` 和 `meta:` 前缀的数据
3. 通过 channel 传递数据
4. 使用 Gin 的 Stream 方法逐块发送

### zhipuHandler(...)
非流式响应处理器：
1. 读取完整响应体
2. 反序列化为 `ZhipuResponse`
3. 检查 `Success` 字段
4. 转换为 OpenAI 格式并写回客户端

## 5. 关键逻辑分析

1. **JWT 缓存**: 使用 `sync.Map` 缓存 Token，有效期 24 小时。Token 过期后自动重新生成。

2. **System 消息转换**: 智谱旧版 API 不支持 system 角色，将其转换为 user+assistant 对话对。

3. **双前缀 SSE**: 智谱的 SSE 流使用两种前缀：`data:` 传递内容，`meta:` 传递元数据（如 Usage）。

4. **Content 清理**: 智谱响应的 Content 可能包含多余引号，需要 `strings.Trim` 清理。

5. **Token 生成**: JWT 使用 HS256 算法，头部包含自定义的 `sign_type: SIGN` 字段。

## 6. 关联文件
- `zhipu/adaptor.go` — 调用处理器函数
- `zhipu/dto.go` — 数据结构定义
- `relay/helper/sse.go` — SSE 流处理工具
