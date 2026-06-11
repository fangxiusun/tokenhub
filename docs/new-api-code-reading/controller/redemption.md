# redemption.go 代码阅读文档

## 1. 全局总结

该文件实现了兑换码（Redemption）的完整 CRUD 管理，包括批量创建、搜索、更新和清理无效兑换码。

## 2. 依赖关系

- `common` — 通用工具、UUID
- `i18n` — 国际化
- `model` — 兑换码模型
- `setting/operation_setting` — 支付合规检查
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

- `GetAllRedemptions` — 分页获取所有兑换码
- `SearchRedemptions` — 搜索兑换码
- `GetRedemption` — 获取单个兑换码
- `AddRedemption` — 批量创建兑换码（最多 100 个）
- `DeleteRedemption` — 删除兑换码
- `UpdateRedemption` — 更新兑换码（支持 status_only 模式）
- `DeleteInvalidRedemption` — 清理无效兑换码

## 5. 关键逻辑分析

- 创建需要支付合规确认
- 名称长度：1-20 个 Unicode 字符
- 数量限制：1-100
- 过期时间不能早于当前时间
- 每个兑换码使用 UUID 作为 Key

## 6. 关联文件

- `model/redemption.go` — 兑换码模型
