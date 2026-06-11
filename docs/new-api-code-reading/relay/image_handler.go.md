# image_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了图像生成请求的处理入口 `ImageHelper`。支持 JSON 和 multipart/form-data 两种请求格式，处理图像请求的转换、发送、响应解析和计费。特别处理了 Replicate 渠道的 201 状态码。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射
- `service`: 计费
- `setting/model_setting`: 全局设置

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ImageHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理图像生成请求
- **特殊处理**:
  - Replicate 渠道返回 201 Created 时视为成功
  - 支持 bytes.Buffer 类型的转换结果（multipart 请求）
  - 通过 OtherRatios 的 "n" 字段处理多图生成的计费
  - 图像请求的 promptTokens 和 totalTokens 至少为 1

## 5. 关键逻辑分析

1. **多图计费**: 通过 `OtherRatios["n"]` 将生成数量纳入计费计算
2. **Replicate 适配**: 201 Created 状态码特殊处理
3. **passthrough 模式**: 支持全局和渠道级别
4. **日志记录**: 记录图像大小、品质、生成数量

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor.ConvertImageRequest 接口
- `dto/image.go`: ImageRequest DTO
