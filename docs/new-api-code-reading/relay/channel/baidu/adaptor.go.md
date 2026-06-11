# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是百度（Baidu）渠道的适配器实现，负责将 OpenAI 格式的请求转换为百度文心一言 API 格式。支持聊天、嵌入等功能，并处理百度特有的访问令牌管理。

## 2. 依赖关系
- 标准库：`errors`, `fmt`, `io`, `net/http`, `strings`
- 内部包：
  - `dto`: 数据传输对象
  - `relay/channel`: 渠道通用工具
  - `relaycommon`: 中继通用配置
  - `relay/constant`: 中继常量
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
### 结构体
- `Adaptor`: 百度渠道适配器结构体（空结构体）。

## 4. 函数详解
### 核心函数
1. **`GetRequestURL`**: 根据模型名称构建请求 URL，支持多种模型端点映射。
2. **`SetupRequestHeader`**: 设置请求头，包括授权信息。
3. **`ConvertOpenAIRequest`**: 转换 OpenAI 请求为百度格式。
4. **`ConvertEmbeddingRequest`**: 转换嵌入请求为百度格式。
5. **`DoRequest`**: 执行 API 请求。
6. **`DoResponse`**: 处理响应，根据流式/非流式和嵌入模式调用不同处理器。

## 5. 关键逻辑分析
- **模型端点映射**：`GetRequestURL` 根据模型名称映射到不同的 API 端点，如 `ERNIE-4.0` 映射到 `completions_pro`。
- **访问令牌管理**：通过 `getBaiduAccessToken` 函数管理百度 API 的访问令牌，支持缓存和自动刷新。
- **URL 构建**：请求 URL 包含访问令牌参数，用于认证。
- **多模式支持**：支持聊天、嵌入和流式响应。

## 6. 关联文件
- `baidu/constants.go`: 定义模型列表和渠道名称。
- `baidu/dto.go`: 百度特定的数据传输对象。
- `baidu/relay-baidu.go`: 请求转换和响应处理逻辑。