# image_wan.go 代码阅读文档

## 1. 全局总结
该文件专门处理阿里云 Wan 系列图像模型的图像编辑请求转换。Wan 模型是阿里云的新一代图像生成模型，具有特定的输入格式和参数要求。

## 2. 依赖关系
- 标准库：`fmt`, `strings`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `relaycommon`: 中继通用配置
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/samber/lo`: 泛型工具库

## 3. 类型定义
该文件没有定义新的类型，主要使用 `ali/dto.go` 中定义的 `WanImageInput` 和 `WanImageParameters`。

## 4. 函数详解
1. **`oaiFormEdit2WanxImageEdit`**: 将表单图像编辑请求转换为 Wan 模型格式，从请求体和表单中提取参数。
2. **`isOldWanModel`**: 判断是否为旧版 Wan 模型（不包含 `wan2.6` 或 `wan2.7`）。
3. **`isWanModel`**: 判断是否为 Wan 模型（包含 `wan` 字符串）。

## 5. 关键逻辑分析
- **模型版本区分**：通过 `isOldWanModel` 和 `isWanModel` 区分不同版本的 Wan 模型，使用不同的 API 端点。
- **输入格式**：Wan 模型使用 `WanImageInput` 结构体，包含提示词、图像数组和反向提示词。
- **价格计算**：在请求转换时计算图像数量比例，用于计费。
- **表单数据处理**：从表单中提取图像 Base64 数据，支持多个图像输入。

## 6. 关联文件
- `ali/adaptor.go`: 调用 `oaiFormEdit2WanxImageEdit` 处理 Wan 模型图像编辑请求。
- `ali/dto.go`: 定义 `WanImageInput` 和 `WanImageParameters` 结构体。
- `ali/image.go`: 包含 `getImageBase64sFromForm` 辅助函数。