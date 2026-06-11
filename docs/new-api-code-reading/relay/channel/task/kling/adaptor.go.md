# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了快影（Kling）视频生成任务的适配器（TaskAdaptor）。快影 API 使用 JWT 认证，支持文生视频和图生视频两种模式，根据 action 类型自动选择不同的 API 端点。适配器还处理了动态动作切换（在 BuildRequestBody 中根据是否有图片决定最终 action）和 token 用量解析。

## 2. 依赖关系

**标准库：**
- `bytes` — 字节缓冲
- `fmt` — 格式化
- `io` — IO 操作
- `math` — 数学函数（Ceil）
- `net/http` — HTTP 请求
- `strconv` — 字符串转换
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
- `github.com/golang-jwt/jwt/v5` — JWT 令牌生成
- `github.com/pkg/errors` — 错误包装
- `github.com/samber/lo` — 泛型工具（Ternary）

## 3. 类型定义

### 请求/响应结构体

| 类型名 | 说明 |
|--------|------|
| `TrajectoryPoint` | 轨迹点，包含 x、y 坐标 |
| `DynamicMask` | 动态遮罩，包含遮罩数据和轨迹点数组 |
| `CameraConfig` | 相机配置，包含水平、垂直、平移、倾斜、旋转、缩放参数 |
| `CameraControl` | 相机控制，包含类型和配置 |
| `requestPayload` | 快影请求体，包含 prompt、图片、负面提示词、模式、时长、宽高比、模型名、CFG 缩放、遮罩、相机控制等 |
| `responsePayload` | 任务查询响应，包含状态码、消息、任务数据（状态、结果视频、时间戳、扣费信息） |

### 适配器结构体

| 类型名 | 说明 |
|--------|------|
| `TaskAdaptor` | 核心适配器，嵌入 `taskcommon.BaseBilling` |

## 4. 函数详解

### 适配器接口方法

| 函数签名 | 说明 |
|----------|------|
| `Init(info)` | 初始化适配器 |
| `ValidateRequestAndSetAction(c, info)` | 验证请求 |
| `BuildRequestURL(info)` | 根据 action 选择端点：image2video 或 text2video |
| `BuildRequestHeader(c, req, info)` | 生成 JWT Token 并设置 Authorization 头 |
| `BuildRequestBody(c, info)` | 转换请求，根据图片有无动态切换 action |
| `DoRequest(c, info, requestBody)` | 从 context 读取动态 action 并执行请求 |
| `DoResponse(c, resp, info)` | 解析响应，检查 code == 0 |
| `FetchTask(baseUrl, key, body, proxy)` | 轮询任务状态，根据 action 选择端点 |
| `GetModelList()` | 返回模型列表 |
| `GetChannelName()` | 返回 "kling" |
| `ParseTaskResult(respBody)` | 解析任务结果，提取视频 URL 和 token 用量 |
| `ConvertToOpenAIVideo(originTask)` | 转换为 OpenAI 格式，包含视频时长 |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `convertToRequestPayload(req, info)` | 转换请求为快影格式 |
| `getAspectRatio(size)` | 将尺寸映射为宽高比 |
| `createJWTToken()` | 使用默认 apiKey 创建 JWT |
| `createJWTTokenWithKey(apiKey)` | 使用指定 apiKey 创建 JWT（30 分钟有效期） |
| `isNewAPIRelay(apiKey)` | 判断是否为 new-api 中继模式 |

## 5. 关键逻辑分析

### JWT 认证
- 使用 HS256 签名算法
- Claims 包含：iss（accessKey）、exp（当前时间+30分钟）、nbf（当前时间-5秒）
- new-api 中继模式（sk- 前缀）直接使用 apiKey 作为 Bearer Token，跳过 JWT 生成

### 动态 Action 切换
1. `BuildRequestURL` 根据初始 action 选择端点
2. `BuildRequestBody` 检查是否有图片，无图片时将 action 改为 TextGenerate 并存入 context
3. `DoRequest` 从 context 读取可能已更新的 action

### 端点选择
- 文生视频：`/v1/videos/text2video`
- 图生视频：`/v1/videos/image2video`
- new-api 中继模式：路径前加 `/kling` 前缀

### 宽高比映射
- 1024x1024 / 512x512 → 1:1
- 1280x720 / 1920x1080 → 16:9
- 720x1280 / 1080x1920 → 9:16
- 其他 → 1:1（默认）

### Token 用量解析
`ParseTaskResult` 从 `FinalUnitDeduction` 字段解析 token 扣费数量，使用 `math.Ceil` 向上取整。

### 状态映射
- submitted → Submitted
- processing → InProgress
- succeed → Success，提取第一个视频的 URL
- failed → Failure

## 6. 关联文件

- `relay/channel/task/taskcommon/` — 任务基类、元数据工具
- `relay/common/relay.go` — TaskSubmitReq 等通用类型
- `dto/video.go` — OpenAIVideo 响应 DTO
