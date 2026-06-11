# token.go 代码阅读文档

## 1. 全局总结

该文件实现了 Token（API 密钥）的完整 CRUD 管理。支持密钥掩码显示、批量操作、IP 限制、模型限制、额度管理等功能。

## 2. 依赖关系

- `common` — 通用工具、密钥生成
- `i18n` — 国际化
- `model` — Token 模型
- `setting/operation_setting` — Token 数量限制
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `TokenBatch` | 批量操作请求 |

## 4. 函数详解

### 查询
- `GetAllTokens` — 分页获取用户 Token 列表
- `SearchTokens` — 搜索 Token
- `GetToken` — 获取单个 Token
- `GetTokenKey` — 获取 Token 完整密钥
- `GetTokenStatus` — 获取 Token 状态（兼容 OpenAI 格式）
- `GetTokenUsage` — 通过 Bearer Token 获取使用量

### CRUD
- `AddToken` — 创建 Token（检查数量限制和额度范围）
- `DeleteToken` — 删除 Token
- `UpdateToken` — 更新 Token（支持 status_only 模式）
- `DeleteTokenBatch` — 批量删除 Token
- `GetTokenKeysBatch` — 批量获取 Token 密钥（最多 100 个）

## 5. 关键逻辑分析

- 所有查询返回掩码后的密钥（`GetMaskedKey`）
- Token 名称最长 50 字符
- 额度上限：10 亿 * QuotaPerUnit
- Token 数量上限通过 `GetMaxUserTokens()` 配置
- 过期或耗尽的 Token 不能重新启用
- 支持跨组重试（`CrossGroupRetry`）

## 6. 关联文件

- `model/token.go` — Token 模型
- `setting/operation_setting/` — Token 配置
