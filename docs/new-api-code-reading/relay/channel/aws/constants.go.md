# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了 AWS Bedrock 渠道的常量和配置，包括模型 ID 映射、跨区域配置和渠道名称。是 AWS 渠道适配器的基础配置文件。

## 2. 依赖关系
- 标准库：`strings`

## 3. 类型定义
### 变量
- `awsModelIDMap`: 模型 ID 映射表，将用户请求的模型名转换为 AWS Bedrock 模型 ID。
- `awsModelCanCrossRegionMap`: 跨区域支持配置，定义每个模型支持的区域。
- `awsRegionCrossModelPrefixMap`: 区域前缀映射表，用于构建跨区域模型 ID。
- `ChannelName`: 渠道名称，值为 `"aws"`。

## 4. 函数详解
1. **`isNovaModel`**: 判断模型 ID 是否为 Nova 模型（包含 `"nova-"` 字符串）。

## 5. 关键逻辑分析
- **模型映射**：`awsModelIDMap` 包含 32 个模型映射，涵盖 Claude 3/3.5/3.7/4 系列和 Nova 系列。
- **跨区域支持**：`awsModelCanCrossRegionMap` 定义每个模型支持的区域（`us`, `eu`, `ap`）。
- **区域前缀**：`awsRegionCrossModelPrefixMap` 将区域代码映射为前缀（如 `ap` 映射为 `apac`）。
- **Nova 模型识别**：通过字符串包含检查识别 Nova 模型。

## 6. 关联文件
- `aws/adaptor.go`: 使用这些常量进行模型 ID 映射和区域配置。
- `aws/relay-aws.go`: 使用这些常量进行 AWS 请求构建。