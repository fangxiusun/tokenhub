# relay-tencent.go 代码阅读文档

## 1. 全局总结
腾讯混元渠道的核心中继逻辑文件。包含请求格式转换、响应格式转换、流式/非流式处理器、API Key 解析、TC3-HMAC-SHA256 签名算法等关键功能。

## 2. 依赖关系
- **标准库**: bufio, crypto/hmac, crypto/sha256, encoding/hex, encoding/json, errors, fmt, io, net/http, strconv, strings, time
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理、日志、时间戳
  - `github.com/QuantumNous/new-api/constant` — 常量定义（FinishReasonStop）
  - `github.com/QuantumNous/new-api/dto` — OpenAI 格式 DTO
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/helper` — 流式扫描器、SSE 头设置
  - `github.com/QuantumNous/new-api/service` — 响应体关闭、使用量计算
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无新增类型定义，使用 `tencent/dto.go` 中定义的类型。

## 4. 函数详解

### requestOpenAI2Tencent(a *Adaptor, request dto.GeneralOpenAIRequest) *TencentChatRequest
将 OpenAI 格式请求转换为腾讯混元格式：
- 遍历消息列表，转换为 `TencentMessage`（提取纯文本内容）
- 复制 Stream、Model、TopP、Temperature 参数

### responseTencent2OpenAI(response *TencentChatResponse) *dto.OpenAITextResponse
将腾讯混元同步响应转换为 OpenAI 格式：
- 设置 Object 为 "chat.completion"
- 提取第一个 Choice 的 Message 内容
- 复制 Usage 信息

### streamResponseTencent2OpenAI(TencentResponse *TencentChatResponse) *dto.ChatCompletionsStreamResponse
将腾讯混元流式响应转换为 OpenAI 流式格式：
- Object 设为 "chat.completion.chunk"
- 从 Delta 字段提取增量内容
- 当 FinishReason 为 "stop" 时设置结束标志

### tencentStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
流式响应处理器：
1. 使用 `helper.NewStreamScanner` 读取 SSE 流
2. 逐行解析 `data:` 前缀的数据
3. 反序列化为 `TencentChatResponse`，转换为 OpenAI 格式
4. 通过 `helper.ObjectData` 写入客户端
5. 最终计算使用量

### tencentHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
非流式响应处理器：
1. 读取完整响应体
2. 反序列化为 `TencentChatResponseSB`
3. 检查错误码
4. 转换为 OpenAI 格式并写回客户端

### parseTencentConfig(config string) (appId int64, secretId string, secretKey string, err error)
解析腾讯云 API 配置，格式为 `appId|secretId|secretKey`，用 `|` 分隔。

### sha256hex(s string) string
计算字符串的 SHA-256 哈希值，返回十六进制字符串。

### hmacSha256(s, key string) string
使用 HMAC-SHA256 计算签名。

### getTencentSign(req TencentChatRequest, adaptor *Adaptor, secId, secKey string) string
TC3-HMAC-SHA256 签名算法实现：
1. 构建规范化请求字符串（HTTP 方法、URI、查询字符串、头部、载荷哈希）
2. 构建签名字符串（算法、时间戳、凭证范围、规范化请求哈希）
3. 使用分层 HMAC 派生签名密钥（日期 → 服务 → 请求）
4. 生成最终 Authorization 头

## 5. 关键逻辑分析

1. **TC3 签名流程**: 腾讯云 API 的签名分为三步：构建规范化请求 → 构建待签字符串 → 分层 HMAC 签名。签名密钥通过 `TC3` + SecretKey → 日期 → 服务名 → 请求类型逐级派生。

2. **SSE 流解析**: 流式处理器使用 `bufio.ScanLines` 逐行读取，跳过长度不足或非 `data:` 前缀的行。支持错误恢复（解析失败时 continue 而非中断）。

3. **API Key 格式**: 腾讯渠道使用三段式 API Key：`appId|secretId|secretKey`，需要在请求前解析。

4. **响应包装**: 非流式响应使用 `TencentChatResponseSB` 包装（外层有 `Response` 字段），而流式响应直接是 `TencentChatResponse`。

## 6. 关联文件
- `tencent/adaptor.go` — 调用本文件的转换函数
- `tencent/dto.go` — 数据结构定义
- `relay/helper/stream_scanner.go` — 流式扫描工具
- `relay/helper/sse.go` — SSE 响应处理工具
- `service/response.go` — 使用量计算和响应体处理
