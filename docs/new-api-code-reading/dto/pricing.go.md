# pricing.go 代码阅读文档

## 1. 全局摘要

该文件定义了不同 AI 服务提供商的模型信息结构体，包括 OpenAI、Anthropic、Gemini 的模型数据结构，用于模型列表和定价信息获取。

## 2. 依赖

- **外部包**：
  - `github.com/QuantumNous/new-api/constant`：常量定义（`EndpointType`）

## 3. 类型定义

### OpenAIModels 结构体
OpenAI 模型信息：
- `Id` (string)：模型 ID
- `Object` (string)：对象类型
- `Created` (int)：创建时间
- `OwnedBy` (string)：所有者
- `SupportedEndpointTypes` ([]constant.EndpointType)：支持的端点类型

### AnthropicModel 结构体
Anthropic 模型信息：
- `ID` (string)：模型 ID
- `CreatedAt` (string)：创建时间
- `DisplayName` (string)：显示名称
- `Type` (string)：模型类型

### GeminiModel 结构体
Gemini 模型信息：
- `Name` (interface{})：模型名称
- `BaseModelId` (interface{})：基础模型 ID
- `Version` (interface{})：版本
- `DisplayName` (interface{})：显示名称
- `Description` (interface{})：描述
- `InputTokenLimit` (interface{})：输入 token 限制
- `OutputTokenLimit` (interface{})：输出 token 限制
- `SupportedGenerationMethods` ([]interface{})：支持的生成方法
- `Thinking` (interface{})：思考能力
- `Temperature` (interface{})：温度
- `MaxTemperature` (interface{})：最大温度
- `TopP` (interface{})：Top-P 采样
- `TopK` (interface{})：Top-K 采样

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **多提供商支持**：为不同 AI 服务提供商定义专用的模型结构体。

2. **灵活类型**：Gemini 模型使用 `interface{}` 类型，适应动态字段类型。

3. **端点类型支持**：OpenAI 模型包含支持的端点类型信息。

## 6. 相关文件

- `constant/endpoint.go`：端点类型常量定义
- `controller/pricing.go`：定价控制器
- `model/pricing.go`：定价数据模型