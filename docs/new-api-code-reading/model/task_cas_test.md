# task_cas_test.go 代码阅读文档

## 1. 全局总结

`task_cas_test.go` 是 `task.go` 中 CAS（Compare-And-Swap）更新机制的测试文件。它包含两部分测试：纯逻辑的快照相等性测试（无数据库依赖）和基于 SQLite 内存数据库的 CAS 更新集成测试。集成测试验证了 CAS 在单次更新、状态不匹配、以及多 goroutine 并发竞争场景下的正确性。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `encoding/json` | JSON 数据构造 |
| `os` | `os.Exit` 用于测试入口 |
| `sync` | `sync.WaitGroup` 等待并发 goroutine 完成 |
| `testing` | Go 测试框架 |
| `time` | 时间戳设置 |
| `github.com/QuantumNous/new-api/common` | 全局配置设置（UsingSQLite 等） |
| `github.com/glebarez/sqlite` | SQLite 驱动（纯 Go 实现） |
| `github.com/stretchr/testify/assert` | 测试断言 |
| `github.com/stretchr/testify/require` | 必须成功的断言 |
| `gorm.io/gorm` | ORM 框架 |

## 3. 类型定义

无额外类型定义。测试使用 `model` 包内已有的类型。

## 4. 函数详解

### TestMain(m *testing.M)
测试入口，初始化测试环境：
1. 打开 SQLite 内存数据库
2. 设置全局 `DB` 和 `LOG_DB`
3. 配置 `common` 包的全局标志（UsingSQLite、RedisEnabled 等）
4. 调用 `initCol()` 初始化列名
5. 设置最大连接数为 1（确保串行执行）
6. 自动迁移所有测试需要的表

### truncateTables(t *testing.T)
测试辅助函数，在每个测试结束后清空所有表。使用 `t.Cleanup` 注册清理函数。

### insertTask(t *testing.T, task *Task)
测试辅助函数，插入任务并设置时间戳。

### 纯逻辑测试（无数据库）

#### TestSnapshotEqual_Same
验证相同快照的 `Equal` 返回 true。

#### TestSnapshotEqual_DifferentStatus
验证不同状态的快照不相等。

#### TestSnapshotEqual_DifferentProgress
验证不同进度的快照不相等。

#### TestSnapshotEqual_DifferentData
验证不同数据的快照不相等。

#### TestSnapshotEqual_NilVsEmpty
验证 nil 和空 `json.RawMessage` 相等（`bytes.Equal` 语义）。

#### TestSnapshot_Roundtrip
验证 `Snapshot()` 方法能正确提取任务的各个字段。

### CAS 集成测试（需要数据库）

#### TestUpdateWithStatus_Win
验证正常 CAS 更新：任务从 `IN_PROGRESS` 更新为 `SUCCESS`，返回 `(true, nil)`，数据库中状态已变更。

#### TestUpdateWithStatus_Lose
验证 CAS 失败：任务当前状态为 `FAILURE`，但 CAS 期望 `IN_PROGRESS`，返回 `(false, nil)`，数据库中状态未变更。

#### TestUpdateWithStatus_ConcurrentWinner
验证并发 CAS：5 个 goroutine 同时尝试从 `IN_PROGRESS` 更新为 `SUCCESS`，断言恰好只有 1 个 goroutine 成功（赢得 CAS）。

## 5. 关键逻辑分析

**测试隔离**：使用 SQLite 内存数据库，每个测试通过 `truncateTables` 清空数据，确保测试之间互不影响。

**并发正确性验证**：`TestUpdateWithStatus_ConcurrentWinner` 是核心测试，通过 5 个并发 goroutine 模拟竞争场景，验证 CAS 机制保证只有一个进程能成功更新。这直接验证了生产环境中计费状态转换的正确性。

**快照纯逻辑测试**：验证 `taskSnapshot.Equal` 的正确性，这是 CAS 机制中用于判断任务是否发生变化的基础。

**数据库串行化**：设置 `MaxOpenConns(1)` 确保 SQLite 的写操作串行执行，使并发测试更可靠。

## 6. 关联文件

- `model/task.go`：被测试的 `Task`、`taskSnapshot`、`UpdateWithStatus`、`Snapshot` 等
- `model/main.go`：`initCol()` 初始化函数
- `model/user.go`、`model/token.go` 等：测试迁移需要的关联模型
