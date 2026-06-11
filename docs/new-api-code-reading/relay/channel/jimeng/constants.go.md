# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了即梦（Jimeng）渠道的常量，包括渠道名称和支持的模型列表。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ChannelName` 常量
```go
const ChannelName = "jimeng"
```
渠道标识名称。

### `ModelList` 变量
```go
var ModelList = []string{
    "jimeng_high_aes_general_v21_L",
}
```
即梦渠道支持的模型列表，目前仅包含一个模型 `jimeng_high_aes_general_v21_L`，这是即梦的高质量通用图像生成模型。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
即梦渠道目前仅支持单一模型，专注于图像生成能力。`req_key` 字段在 `adaptor.go` 的 `ConvertImageRequest` 中直接使用模型名称。

## 6. 关联文件
- `adaptor.go` — 使用 `ModelList` 和 `ChannelName`
