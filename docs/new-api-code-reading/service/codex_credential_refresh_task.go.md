# codex_credential_refresh_task.go 代码阅读文档

## 1. 全局总结

该文件实现了 Codex 凭证的自动刷新定时任务。每 10 分钟检查一次即将过期的 Codex 通道，自动刷新其 OAuth 令牌。仅在主节点上运行。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 主节点判断 |
| `constant` | 通道类型常量 |
| `logger` | 日志记录 |
| `model` | 数据库查询 |
| `gopool` | 异步任务 |

## 3. 类型定义

无自定义类型，使用常量：
- `codexCredentialRefreshTickInterval = 10 * time.Minute` — 检查间隔
- `codexCredentialRefreshThreshold = 24 * time.Hour` — 过期阈值
- `codexCredentialRefreshBatchSize = 200` — 批量大小
- `codexCredentialRefreshTimeout = 15 * time.Second` — 单次刷新超时

## 4. 函数详解

### `StartCodexCredentialAutoRefreshTask()`
- 使用 `sync.Once` 确保只启动一次
- 仅主节点运行
- 启动 ticker 定时执行

### `runCodexCredentialAutoRefreshOnce()`
- 使用 CAS 防止并发执行
- 分批查询 Codex 类型通道（enabled 或 auto_disabled 状态）
- 跳过多密钥通道
- 检查过期时间，距过期超过 24 小时则跳过
- 逐个刷新凭证
- 刷新完成后重置通道缓存和代理客户端缓存

### `shouldAutoRefreshCodexChannelStatus(status int) bool`
判断通道状态是否需要自动刷新

## 5. 关键逻辑分析

1. **CAS 防重入**：使用 `atomic.Bool` 确保同一时间只有一个刷新任务在运行
2. **批量处理**：每次最多处理 200 个通道
3. **过期阈值**：仅刷新距过期不足 24 小时的通道
4. **缓存刷新**：刷新后批量重置缓存（带 recover 保护）

## 6. 关联文件

- `codex_credential_refresh.go` — 单个通道凭证刷新
- `codex_oauth.go` — OAuth 令牌刷新实现
