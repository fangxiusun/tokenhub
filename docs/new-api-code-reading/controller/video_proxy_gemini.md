# video_proxy_gemini.go 代码阅读文档

## 1. 全局总结

该文件实现了 Gemini 视频 URL 的获取和代理功能。从任务数据中提取视频 URL，或通过 API 获取视频下载链接。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — 常量
- `model` — 任务和渠道模型
- `relay` — 任务适配器

## 3. 类型定义

无。

## 4. 函数详解

### `getGeminiVideoURL(channel, task, apiKey) (string, error)`
获取 Gemini 视频 URL。优先从任务数据中提取，否则通过 API 获取。

### `extractGeminiVideoURLFromTaskData(task) string`
从任务数据中提取视频 URL。

### `ensureAPIKey(url, apiKey) string`
确保 URL 包含 API Key 参数。

## 5. 关键逻辑分析

- 优先使用缓存的任务数据中的 URL
- URL 需要附加 API Key 用于认证
- 支持多种视频 URL 格式

## 6. 关联文件

- `controller/video_proxy.go` — 视频代理
- `relay/channel/gemini/` — Gemini 适配器
