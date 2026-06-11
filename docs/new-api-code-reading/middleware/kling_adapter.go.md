# kling_adapter.go 代码阅读文档

## 1. 全局总结

该文件实现了 Kling AI API 请求格式转换中间件 `KlingRequestConvert`，将 Kling 原生 API 请求格式转换为统一的视频生成请求格式。与 `jimeng_adapter.go` 类似，但更简洁，仅支持提交任务模式。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `bytes` | 构建请求体缓冲区 |
| `encoding/json` | JSON 序列化 |
| `io` | 请求体替换 |
| `github.com/QuantumNous/new-api/common` | 请求体解析、常量定义 |
| `github.com/QuantumNous/new-api/constant` | 任务动作常量 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `KlingRequestConvert() func(c *gin.Context)`

- **功能**：创建 Kling 请求转换中间件。
- **执行流程**：
  1. 解析原始请求体
  2. 支持 `model_name` 和 `model` 两种字段名获取模型名称
  3. 提取 `prompt` 字段
  4. 构建统一请求格式：`{model, prompt, metadata: 原始请求}`
  5. 替换请求体并将路径改为 `/v1/video/generations`
  6. 如果没有 `image` 字段，设置 action 为文本生成

## 5. 关键逻辑分析

- **兼容性设计**：同时支持 `model_name` 和 `model` 两种字段名，提高对不同 Kling API 版本的兼容性。
- **静默降级**：请求解析失败时直接调用 `c.Next()` 继续处理，不中断请求（与 Jimeng 适配器的错误处理策略不同）。
- **请求体重写**：与 Jimeng 适配器相同，将原始请求体作为 `metadata` 保留。
- **任务类型判断**：根据是否包含 `image` 字段判断文本生成或图像生成。

## 6. 关联文件

- `middleware/jimeng_adapter.go` — 类似的 Jimeng API 适配器
- `common/request.go` — `UnmarshalBodyReusable` 和 `KeyRequestBody` 定义
