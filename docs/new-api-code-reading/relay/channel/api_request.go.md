# api_request.go 代码阅读文档

## 1. 全局总结

本文件实现了 Relay 模块的 HTTP 请求发送核心逻辑，包括标准 API 请求、表单请求、WebSocket 请求和任务请求。还实现了请求头覆盖（Header Override）系统，支持通配符、正则匹配和客户端请求头透传。

## 2. 依赖关系

- `relay/common`: RelayInfo
- `relay/helper`: SSE ping、流式头设置
- `service`: HTTP 客户端、代理
- `operation_setting`: Ping 间隔配置
- `gorilla/websocket`: WebSocket 客户端

## 3. 类型定义

本文件无自定义类型定义，但定义了重要的常量和变量：
- `passthroughSkipHeaderNamesLower`: 不应透传的 header 名称集合（hop-by-hop、认证、WebSocket 握手头等）
- `headerPassthroughRegexCache`: 正则表达式缓存

## 4. 函数详解

### 核心请求函数

#### `DoApiRequest(a, c, info, requestBody) (*http.Response, error)`
- 标准 API 请求入口，处理 URL 构建、请求头设置、Header Override

#### `DoFormRequest(a, c, info, requestBody) (*http.Response, error)`
- 表单请求入口，额外设置 Content-Type

#### `DoWssRequest(a, c, info, requestBody) (*websocket.Conn, error)`
- WebSocket 请求入口，建立 WS 连接

#### `DoTaskApiRequest(a, c, info, requestBody) (*http.Response, error)`
- 任务请求入口，使用 TaskAdaptor 接口

### Header Override 系统

#### `processHeaderOverride(info, c) (map[string]string, error)`
- 处理请求头覆盖规则
- 支持通配符 `*`、正则 `re:` / `regex:`、客户端请求头透传 `{client_header:name}`

#### `applyHeaderOverrideToRequest(req, headerOverride)`
- 将覆盖规则应用到 HTTP 请求

#### `applyHeaderOverridePlaceholders(template, c, apiKey) (string, bool, error)`
- 处理占位符替换：`{api_key}` → 渠道密钥，`{client_header:name}` → 客户端请求头

### Ping 保活

#### `startPingKeepAlive(c, pingInterval) context.CancelFunc`
- 启动 SSE ping 保活协程
- 支持超时保护（120分钟最大持续时间）

#### `sendPingData(c, mutex) error`
- 发送 ping 数据，带超时控制（10秒）

## 5. 关键逻辑分析

1. **Header Override 优先级**: passthrough 规则先应用，显式覆盖后应用（显式优先）
2. **安全跳过**: 不透传认证头、hop-by-hop 头、WebSocket 握手头
3. **Content-Length 设置**: 当 body 是 BodyStorage 类型时，手动设置 Content-Length
4. **Ping 保活**: 在流式请求期间定期发送 SSE ping，防止连接超时
5. **代理支持**: 通过渠道设置的 Proxy 配置创建代理 HTTP 客户端

## 6. 关联文件

- `relay/common/override.go`: 参数覆盖逻辑
- `relay/helper/common.go`: SSE 头设置
- `service/http_client.go`: HTTP 客户端管理
