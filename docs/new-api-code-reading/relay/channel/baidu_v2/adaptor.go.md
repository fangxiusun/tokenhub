# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是百度 V2（Volcengine/火山引擎）渠道的适配器实现，负责将 OpenAI 格式的请求转换为百度 V2 API 格式。支持聊天、嵌入、图像生成、重排序等功能，并实现了 `channel.Adaptor` 接口。

## 2. 依赖关系
- 标准库：`errors`, `fmt`, `io`, `net/http`, `strings`
- 内部包：
  - `dto`: 数据传输对象
  - `relay/channel`: 渠道通用工具
  - `relay/channel/openai`: OpenAI 渠道适配器
  - `relaycommon`: 中继通用配置
  - `relay/constant`: 中继常量
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
### 结构体
- `Adaptor`: 百度 V2 渠道适配器结构体（空结构体）。

## 4. 函数详解
### 核心函数
1. **`ConvertClaudeRequest`**: 委托给 OpenAI 适配器处理 Claude 请求。
2. **`GetRequestURL`**: 根据中继模式构建请求 URL，支持多种端点（聊天、嵌入、图像生成等）。
3. **`SetupRequestHeader`**: 设置请求头，处理 API Key 和 AppID。
4. **`ConvertOpenAIRequest`**: 转换 OpenAI 请求，支持搜索功能。
5. **`DoRequest`**: 执行 API 请求。
6. **`DoResponse`**: 委托给 OpenAI 适配器处理响应。

## 5. 关键逻辑分析
- **搜索功能支持**：模型名称以 `-search` 结尾时，自动启用网络搜索功能。
- **API Key 格式**：API Key 格式为 `token|appid`，支持可选的 AppID。
- **URL 构建**：根据不同中继模式构建对应的 API 端点 URL。
- **委托模式**：Claude 请求和响应处理委托给 OpenAI 适配器。

## 6. 关联文件
- `baidu_v2/constants.go`: 定义模型列表和渠道名称。
- `relay/channel/openai/adaptor.go`: OpenAI 适配器，用于 Claude 请求和响应处理。
- `relay/constant/constant.go`: 中继常量定义。