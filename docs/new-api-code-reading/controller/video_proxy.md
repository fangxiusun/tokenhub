# video_proxy.go 代码阅读文档

## 1. 全局总结

该文件实现了视频内容的代理服务，通过后端代理前端请求视频 URL，解决跨域和认证问题。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — 常量
- `logger` — 日志
- `model` — 任务模型
- `service` — HTTP 客户端
- `setting/system_setting` — 服务器地址
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `videoProxyError(c, status, errType, message)`
返回标准化的 OpenAI 风格错误响应。

### 视频代理处理器
代理前端请求视频内容，支持范围请求（Range header）和流式传输。

## 5. 关键逻辑分析

- 支持 HTTP Range 请求（视频播放器需要）
- 使用服务端代理解决跨域问题
- 支持 base64 编码的 URL

## 6. 关联文件

- `controller/video_proxy_gemini.go` — Gemini 视频代理
- `model/task.go` — 任务模型
