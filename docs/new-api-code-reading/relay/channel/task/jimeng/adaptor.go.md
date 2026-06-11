# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了即梦（Jimeng/字节跳动视觉生成平台）视频生成任务的适配器（TaskAdaptor）。即梦 API 使用火山引擎的 CV（计算机视觉）同步转异步接口，支持两种认证方式：传统的 AccessKey/SecretKey HMAC 签名认证和 Bearer Token 认证（通过 new-api 中继）。适配器还处理了即梦视频 3.0 的 ReqKey 转换逻辑。

## 2. 依赖关系

**标准库：**
- `bytes` — 字节缓冲
- `crypto/hmac` — HMAC 签名
- `crypto/sha256` — SHA256 哈希
- `encoding/base64` — Base64 编解码
- `encoding/hex` — 十六进制编码
- `fmt` — 格式化
- `io` — IO 操作
- `net/http` — HTTP 请求
- `net/url` — URL 编码
- `sort` — 排序
- `strings` — 字符串处理
- `time` — 时间处理

**项目内部依赖：**
- `github.com/QuantumNous/new-api/common` — JSON 序列化
- `github.com/QuantumNous/new-api/model` — 任务状态
- `github.com/QuantumNous/new-api/constant` — 任务动作常量
- `github.com/QuantumNous/new-api/dto` — OpenAI 视频 DTO
- `github.com/QuantumNous/new-api/relay/channel` — 通用请求执行
- `github.com/QuantumNous/new-api/relay/channel/task/taskcommon` — 任务基类、元数据工具
- `github.com/QuantumNous/new-api/relay/common` — 中继信息、任务请求
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

**第三方依赖：**
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/pkg/errors` — 错误包装
- `github.com/samber/lo` — 泛型工具（Max）

## 3. 类型定义

### 请求/响应结构体

| 类型名 | 说明 |
|--------|------|
| `requestPayload` | 即梦请求体，包含 req_key、二进制数据、图片 URL、提示词、种子、宽高比、帧数 |
| `responsePayload` | 提交响应，包含状态码、消息、请求 ID、任务 ID |
| `responseTask` | 任务查询响应，包含状态码、数据（视频 URL、状态）、消息、耗时 |

### 常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `MaxFileSize` | 4.7MB | 即梦单文件最大限制 |

### 适配器结构体

| 类型名 | 说明 |
|--------|------|
| `TaskAdaptor` | 核心适配器，包含 accessKey、secretKey（非标准 apiKey 格式） |

## 4. 函数详解

### 适配器接口方法

| 函数签名 | 说明 |
|----------|------|
| `Init(info)` | 初始化，解析 "access_key\|secret_key" 格式的 API 密钥 |
| `ValidateRequestAndSetAction(c, info)` | 验证请求 |
| `BuildRequestURL(info)` | 构建 URL，支持两种路径格式 |
| `BuildRequestHeader(c, req, info)` | 设置请求头，支持 Bearer Token 和 HMAC 签名两种认证 |
| `BuildRequestBody(c, info)` | 处理 multipart 文件上传，转换请求格式 |
| `DoRequest(c, info, requestBody)` | 委托通用请求执行 |
| `DoResponse(c, resp, info)` | 解析响应，检查 code == 10000 表示成功 |
| `FetchTask(baseUrl, key, body, proxy)` | 轮询任务状态，发送 POST 请求 |
| `GetModelList()` | 返回模型列表（仅 `jimeng_vgfm_t2v_l20`） |
| `GetChannelName()` | 返回 "jimeng" |
| `ParseTaskResult(respBody)` | 解析任务结果 |
| `ConvertToOpenAIVideo(originTask)` | 转换为 OpenAI 格式 |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `signRequest(req, accessKey, secretKey) error` | 火山引擎 HMAC-SHA256 签名 |
| `hmacSHA256(key, data) []byte` | HMAC-SHA256 计算 |
| `convertToRequestPayload(req, info)` | 转换请求，处理 ReqKey 映射 |
| `isNewAPIRelay(apiKey) bool` | 判断是否为 new-api 中继模式（sk- 前缀） |

## 5. 关键逻辑分析

### 双认证模式
- **new-api 中继模式**：apiKey 以 "sk-" 开头，使用 Bearer Token 认证，URL 路径带 `/jimeng/` 前缀
- **直连模式**：apiKey 格式为 "access_key|secret_key"，使用火山引擎 HMAC-SHA256 签名认证

### HMAC 签名流程（火山引擎标准）
1. 计算请求体的 SHA256 哈希
2. 构建规范请求（方法、路径、查询参数、头部、哈希）
3. 构建待签名字符串（HMAC-SHA256 + 日期 + 凭证范围 + 规范请求哈希）
4. 派生签名密钥（secretKey → date → region → service → request）
5. 生成签名并设置 Authorization 头

### 即梦视频 3.0 ReqKey 转换
根据图片数量自动转换 ReqKey：
- 无图片：`jimeng_v30` → `jimeng_t2v_v30`（文生视频）
- 1 张图片：`jimeng_v30` → `jimeng_i2v_first_v30`（图生视频）
- 多张图片：`jimeng_v30` → `jimeng_i2v_first_tail_v30`（首尾帧生视频）
- `jimeng_v30_pro` 固定转换为 `jimeng_ti2v_v30_pro`

### Multipart 文件处理
支持 OpenAI SDK 的图片上传方式，从 `input_reference` 字段读取文件，检查大小限制（4.7MB），转换为 base64 编码。单文件为图生视频，多文件为首尾帧生视频。

### 帧数计算
- 5 秒视频：24*5+1 = 121 帧
- 10 秒视频：24*10+1 = 241 帧

## 6. 关联文件

- `relay/channel/task/taskcommon/` — 元数据反序列化工具
- `relay/common/relay.go` — TaskSubmitReq 等通用类型
- `dto/video.go` — OpenAIVideo 响应 DTO
