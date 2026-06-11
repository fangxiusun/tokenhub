# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了百度 V2（Volcengine/火山引擎）渠道的常量和模型列表，是百度 V2 渠道适配器的基础配置文件。文件包含两个主要变量：`ModelList` 和 `ChannelName`，用于标识该渠道支持的模型和渠道名称。

## 2. 依赖关系
该文件没有外部依赖，仅使用 Go 标准库的字符串切片定义。

## 3. 类型定义
### 变量
- `ModelList`: 字符串切片，包含百度 V2 渠道支持的所有模型名称。
- `ChannelName`: 字符串常量，值为 `"volcengine"`，标识渠道名称。

## 4. 函数详解
该文件没有定义任何函数。

## 5. 关键逻辑分析
- `ModelList` 包含了 22 个模型，涵盖 ERNIE 系列、DeepSeek 系列等。
- 模型命名遵循百度 V2 的命名规范，如 `ernie-4.0-8k-latest`、`deepseek-v3` 等。
- `ChannelName` 为 `"volcengine"`（火山引擎），表明这是百度 V2 的云服务品牌。

## 6. 关联文件
- `baidu_v2/adaptor.go`: 百度 V2 渠道适配器，使用这些常量进行请求构建和响应处理。
- `relay/constant/constant.go`: 全局渠道常量定义，可能引用 `ChannelName`。
- `model/channel.go`: 渠道模型，可能使用 `ModelList` 进行模型验证。