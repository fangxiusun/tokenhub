# playground.go 代码阅读文档

## 1. 全局总结

该文件实现了 Playground（测试沙箱）功能，允许已登录用户直接测试 AI 模型而无需创建 Token。

## 2. 依赖关系

- `middleware` — Token 上下文设置
- `model` — 用户缓存、Token 模型
- `relay/common` — RelayInfo 生成
- `types` — 错误类型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `Playground(c *gin.Context)`
Playground 入口。创建临时 Token 对象，设置用户上下文，然后调用 `Relay` 处理请求。

## 5. 关键逻辑分析

- 不支持 Access Token 认证（仅支持 Session 认证）
- 创建临时 Token 包含 userId、名称和组信息
- 名称格式：`playground-{group}`
- 最终调用标准的 `Relay` 函数处理请求

## 6. 关联文件

- `controller/relay.go` — `Relay` 函数
