# relay_utils.go 代码阅读文档

## 1. 全局总结

本文件提供了 Relay 模块的通用工具函数，包括请求 URL 构建、API 版本获取、任务请求验证和解析等。

## 2. 依赖关系

- `common`: 通用工具
- `constant`: 常量
- `dto`: 任务 DTO
- `gin`: HTTP 框架

## 3. 类型定义

### 接口
- `HasPrompt`: 包含 GetPrompt() 方法的接口
- `HasImage`: 包含 HasImage() 方法的接口

## 4. 函数详解

### `GetFullRequestURL(baseURL, requestURL, channelType) string`
- 构建完整的上游请求 URL
- 特殊处理 Cloudflare AI Gateway 的路径格式

### `GetAPIVersion(c) string`
- 从查询参数或上下文获取 API 版本

### `ValidateMultipartDirect(c, info) *dto.TaskError`
- 验证 Sora 等视频任务的直接提交请求
- 支持 JSON body 格式
- 处理 model、size、duration 等参数验证

### `ValidateBasicTaskRequest(c, info, action) *dto.TaskError`
- 验证基本任务请求
- 支持 multipart/form-data 和 JSON 两种格式
- 自动兼容单图上传（image → images）

## 5. 关键逻辑分析

1. **Cloudflare 适配**: OpenAI 渠道去掉 `/v1` 前缀，Azure 渠道去掉 `/openai/deployments` 前缀
2. **Sora 验证**: sora-2 仅支持 720x1280/1280x720，sora-2-pro 额外支持 1792x1024/1024x1792
3. **Metadata 兼容**: 支持 Metadata 为字符串或对象格式

## 6. 关联文件

- `relay/relay_task.go`: 任务提交逻辑
- `dto/task.go`: 任务 DTO
