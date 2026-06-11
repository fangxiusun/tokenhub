# image.go 代码阅读文档

## 1. 全局总结
该文件实现了阿里云（Ali）渠道的图像请求和响应处理功能，包括将 OpenAI 格式的图像请求转换为阿里云格式、处理图像编辑请求、异步任务轮询以及响应转换。是阿里云图像生成和编辑功能的核心实现。

## 2. 依赖关系
- 标准库：`encoding/base64`, `errors`, `fmt`, `io`, `mime/multipart`, `net/http`, `strings`, `time`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `github.com/QuantumNous/new-api/logger`: 日志记录
  - `relaycommon`: 中继通用配置
  - `service`: 业务逻辑服务
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/samber/lo`: 泛型工具库

## 3. 类型定义
该文件没有定义新的类型，主要使用 `ali/dto.go` 中定义的类型。

## 4. 函数详解
### 请求转换函数
1. **`oaiImage2AliImageRequest`**: 将 OpenAI 图像请求转换为阿里云图像请求，支持同步和异步模式。
2. **`getImageBase64sFromForm`**: 从表单数据中提取图像的 Base64 编码。
3. **`oaiFormEdit2AliImageEdit`**: 将表单图像编辑请求转换为阿里云格式。

### 任务管理函数
4. **`updateTask`**: 更新异步任务状态。
5. **`asyncTaskWait`**: 等待异步任务完成，支持超时控制。

### 响应处理函数
6. **`responseAli2OpenAIImage`**: 将阿里云图像响应转换为 OpenAI 格式。
7. **`aliImageHandler`**: 处理图像响应，支持同步和异步任务。

## 5. 关键逻辑分析
- **同步/异步处理**：根据模型类型（`isSync`）选择不同的请求格式和处理流程。
- **表单数据处理**：支持从 `multipart/form-data` 中提取图像数据，兼容多种表单字段名。
- **任务轮询机制**：异步任务使用轮询方式检查状态，最大轮询 20 次，每次间隔 10 秒。
- **价格计算**：在请求转换时计算价格比例（如 `prompt_extend`、`n` 等），用于后续计费。
- **错误处理**：在关键步骤返回明确的错误信息，便于调试和错误追踪。

## 6. 关联文件
- `ali/adaptor.go`: 调用这些函数处理图像请求和响应。
- `ali/dto.go`: 定义图像相关的数据传输对象。
- `ali/image_wan.go`: Wan 模型图像编辑处理。
- `service/image.go`: 图像处理服务，可能包含 `GetImageFromUrl` 等辅助函数。
- `relay/common/relay_info.go`: 中继信息结构体，包含价格计算相关字段。