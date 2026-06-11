# client.go 代码阅读文档

## 1. 全局总结
该文件实现了 IO.NET API 客户端，提供 HTTP 请求执行、客户端构造和请求参数构建功能。支持企业版和公共版 API 端点，包含默认超时和错误处理。

## 2. 依赖关系
- **bytes**: 缓冲区操作。
- **encoding/json**: JSON 序列化。
- **fmt**: 格式化错误信息。
- **net/http**: HTTP 客户端。
- **net/url**: URL 参数构建。
- **strconv**: 类型转换。
- **time**: 超时时间处理。

## 3. 类型定义
### DefaultHTTPClient
结构体，实现 HTTPClient 接口，封装标准 http.Client。

## 4. 函数详解
### NewDefaultHTTPClient
构造函数，创建带超时的默认 HTTP 客户端。

### Do
执行 HTTP 请求，设置请求头，返回响应。处理响应体读取和头转换。

### NewEnterpriseClient
创建企业版 API 客户端，使用 DefaultEnterpriseBaseURL。

### NewClient
创建公共版 API 客户端，使用 DefaultBaseURL。

### NewClientWithConfig
创建自定义配置的客户端，支持自定义 baseURL 和 HTTP 客户端。

### makeRequest
内部方法，执行 HTTP 请求，处理 JSON 序列化、错误响应解析（支持 {"detail": "message"} 格式）。

### buildQueryParams
构建 GET 请求的查询参数，支持多种类型（字符串、整数、浮点数、布尔值、时间、数组等）。

## 5. 关键逻辑分析
- **错误处理**: API 错误响应解析为 APIError 结构体，支持详细错误信息。
- **参数构建**: buildQueryParams 支持多种类型，自动跳过零值和空值。
- **超时控制**: 默认客户端超时 30 秒。

## 6. 关联文件
- **types.go**: 定义 Client、HTTPClient、HTTPRequest、HTTPResponse 等类型。
- **container.go、deployment.go、hardware.go**: 使用 Client 的 makeRequest 方法执行 API 调用。
- **jsonutil.go**: 提供 JSON 解析工具函数。