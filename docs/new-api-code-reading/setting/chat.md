# chat.go 代码阅读文档

## 1. 全局总结

该文件管理聊天客户端配置列表，支持多种第三方聊天客户端（如 Cherry Studio、Lobe Chat、OpenCat 等）的快速连接配置。

## 2. 依赖关系

- `encoding/json` — JSON 序列化/反序列化
- `github.com/QuantumNous/new-api/common` — 系统日志

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `Chats` | `[]map[string]string` | 聊天客户端配置列表，每个元素包含客户端名称和 URL 模板 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `UpdateChatsByJsonString` | `func UpdateChatsByJsonString(jsonString string) error` | 从 JSON 字符串更新聊天客户端配置 |
| `Chats2JsonString` | `func Chats2JsonString() string` | 将聊天客户端配置序列化为 JSON 字符串 |

## 5. 关键逻辑分析

- 每个聊天客户端配置是一个 map，key 为客户端名称，value 为 URL 模板
- URL 模板支持变量占位符：`{key}`、`{address}`、`{cherryConfig}` 等
- 该文件直接使用 `encoding/json` 而非 `common` 包的 JSON 封装

## 6. 关联文件

- `common/json.go` — JSON 操作封装
- `controller/option.go` — 管理界面配置接口
