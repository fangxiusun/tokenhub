# channel_affinity.go 代码阅读文档

## 1. 全局总结

该文件实现了通道亲和性（Channel Affinity）系统，允许根据规则将特定请求绑定到特定通道。支持多种匹配条件（模型正则、路径正则、User-Agent、请求头、JSON 字段等），通过混合缓存（内存+Redis）实现高性能的通道绑定和使用统计。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `cachex` | 混合缓存（内存+Redis） |
| `hot` | 内存 LRU 缓存 |
| `gjson` | JSON 字段提取 |
| `operation_setting` | 亲和性规则配置 |
| `common` | SHA1 哈希、Redis 操作 |
| `gin` | HTTP 上下文 |
| `sync` | 并发控制 |

## 3. 类型定义

### `channelAffinityMeta`
亲和性元数据，存储在 gin.Context 中：
- `CacheKey` / `TTLSeconds` — 缓存键和过期时间
- `RuleName` / `SkipRetry` — 规则名称和重试策略
- `ParamTemplate` — 参数覆盖模板
- `KeySourceType/Key/Path` — 键来源信息
- `KeyFingerprint` — 键指纹（SHA1前8位）

### `ChannelAffinityStatsContext`
统计上下文，用于使用量缓存统计

### `ChannelAffinityCacheStats`
缓存统计信息（启用状态、总条目数、按规则分布等）

### `ChannelAffinityUsageCacheStats` / `ChannelAffinityUsageCacheCounters`
使用量缓存统计（命中率、token用量等）

## 4. 核心函数

### `GetPreferredChannelByAffinity(c, modelName, usingGroup) (int, bool)`
- 遍历规则列表，匹配模型正则、路径正则、User-Agent
- 从 KeySource 提取亲和性值（context/header/gjson）
- 查找缓存中的通道 ID
- 将元数据写入 gin.Context

### `RecordChannelAffinity(c, channelID)`
- 请求成功后记录通道亲和性
- 支持 `SwitchOnSuccess` 配置（使用实际成功的通道ID）

### `ApplyChannelAffinityOverrideTemplate(c, paramOverride) (map, bool)`
- 合并规则中的参数覆盖模板到通道配置

### `ObserveChannelAffinityUsageCacheByRelayFormat(c, usage, relayFormat)`
- 记录使用量缓存统计（命中率、token用量）

## 5. 关键逻辑分析

1. **规则优先级**：按配置顺序遍历规则，第一个匹配的规则生效
2. **混合缓存**：内存 LRU + Redis，支持高并发场景
3. **指纹机制**：使用 SHA1 前8位作为指纹，避免存储完整键值
4. **分片锁**：使用 FNV 哈希分片的锁数组减少锁竞争
5. **正则缓存**：使用 sync.Map 缓存编译后的正则表达式

## 6. 关联文件

- `channel_select.go` — 通道选择逻辑
- `operation_setting/channel_affinity.go` — 规则配置
- `pkg/cachex` — 混合缓存实现
