# jimeng_adapter.go 代码阅读文档

## 1. 全局总结

该文件实现了即梦（Jimeng）AI API 请求格式转换中间件 `JimengRequestConvert`，将即梦原生 API 请求格式转换为统一的视频生成请求格式，支持同步请求和异步任务结果查询两种模式。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `bytes` | 构建请求体缓冲区 |
| `encoding/json` | JSON 序列化 |
| `io` | 请求体替换 |
| `net/http` | HTTP 方法和状态码 |
| `github.com/QuantumNous/new-api/common` | 请求体解析、常量定义 |
| `github.com/QuantumNous/new-api/constant` | 任务动作常量 |
| `github.com/QuantumNous/new-api/relay/constant` | Relay 模式常量 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `JimengRequestConvert() func(c *gin.Context)`

- **功能**：创建即梦请求转换中间件。
- **执行流程**：
  1. 从 URL query 获取 `Action` 参数
  2. 解析原始请求体，提取 `req_key`（模型）和 `prompt`
  3. 构建统一请求格式：`{model, prompt, metadata: 原始请求}`
  4. 替换请求体并修改 URL 路径为 `/v1/video/generations`
  5. 如果没有 `image` 字段，设置 action 为文本生成
  6. 对于 `CVSync2AsyncGetResult` action，修改为 GET 请求并设置任务查询路径

## 5. 关键逻辑分析

- **Action 路由**：
  - 其他 Action → `POST /v1/video/generations`（提交任务）
  - `CVSync2AsyncGetResult` → `GET /v1/video/generations/{taskId}`（查询任务结果）
- **请求体重写**：将原始请求体作为 `metadata` 字段保留，同时提取 `model` 和 `prompt` 为顶层字段。
- **任务类型判断**：根据请求中是否包含 `image` 字段判断是文本生成还是图像生成。
- **错误处理**：使用 `abortWithOpenAiMessage` 返回 OpenAI 兼容的错误格式。

## 6. 关联文件

- `middleware/kling_adapter.go` — 类似的 Kling API 适配器
- `relay/constant/relay_mode.go` — `RelayModeVideoFetchByID` 常量定义
- `common/request.go` — `UnmarshalBodyReusable` 和 `KeyRequestBody` 定义
