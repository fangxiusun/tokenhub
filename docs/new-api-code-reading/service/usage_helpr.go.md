# usage_helpr.go 代码阅读文档

## 1. 全局总结

该文件提供 Usage 相关的辅助函数，包括从响应文本生成 Usage 和验证 Usage 有效性。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 上下文操作 |
| `constant` | 上下文键名 |
| `dto` | Usage 结构体 |

## 3. 函数详解

### `ResponseText2Usage(c, responseText, modelName, promptTokens) *dto.Usage`
从响应文本生成 Usage：
- 设置 ContextKeyLocalCountTokens 标志
- 使用估算器计算 completion tokens

### `ValidUsage(usage) bool`
验证 Usage 是否有效（prompt 或 completion 不为 0）

## 4. 关联文件

- `token_estimator.go` — Token 估算
- `dto/usage.go` — Usage 结构体
