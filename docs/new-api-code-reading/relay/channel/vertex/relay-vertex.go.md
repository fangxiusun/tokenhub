# relay-vertex.go 代码阅读文档

## 1. 全局总结
Vertex AI 渠道的辅助逻辑文件，提供模型区域（Region）获取功能。支持从 JSON 配置中按模型名查找区域，或使用默认区域。

## 2. 依赖关系
- **内部包**: `github.com/QuantumNous/new-api/common`

## 3. 类型定义
无类型定义。

## 4. 函数详解

### GetModelRegion(other string, localModelName string) string
获取模型对应的 Google Cloud 区域：
- 如果 `other` 是 JSON 对象字符串：
  - 先尝试按 `localModelName` 查找区域
  - 找不到则使用 `default` 键的值
  - 都没有则返回 `"global"`
- 如果 `other` 不是 JSON，直接返回 `other` 作为区域名

## 5. 关键逻辑分析
- 支持灵活的区域配置：可以为不同模型指定不同区域
- JSON 配置格式示例：`{"model-name": "us-central1", "default": "global"}`
- `common.IsJsonObject` 用于判断字符串是否为 JSON 格式
- `common.StrToMap` 将 JSON 字符串解析为 map

## 6. 关联文件
- `vertex/url_builder.go` — 使用 `GetModelRegion` 获取区域来构建 URL
- `vertex/adaptor.go` — 在 `getRequestUrl` 中调用
