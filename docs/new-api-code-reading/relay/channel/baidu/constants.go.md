# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了百度（Baidu）渠道的常量和模型列表，是百度渠道适配器的基础配置文件。文件包含两个主要变量：`ModelList` 和 `ChannelName`，用于标识该渠道支持的模型和渠道名称。

## 2. 依赖关系
该文件没有外部依赖，仅使用 Go 标准库的字符串切片定义。

## 3. 类型定义
### 变量
- `ModelList`: 字符串切片，包含百度渠道支持的所有模型名称。
- `ChannelName`: 字符串常量，值为 `"baidu"`，标识渠道名称。

## 4. 函数详解
该文件没有定义任何函数。

## 5. 关键逻辑分析
- `ModelList` 包含了 16 个模型，涵盖 ERNIE 系列、BLOOMZ、嵌入模型等。
- 模型命名遵循百度的命名规范，如 `ERNIE-4.0-8K`、`bge-large-zh` 等。
- `ChannelName` 用于在系统中标识百度渠道，便于路由和配置管理。

## 6. 关联文件
- `baidu/adaptor.go`: 百度渠道适配器，使用这些常量进行请求构建和响应处理。
- `relay/constant/constant.go`: 全局渠道常量定义，可能引用 `ChannelName`。
- `model/channel.go`: 渠道模型，可能使用 `ModelList` 进行模型验证。