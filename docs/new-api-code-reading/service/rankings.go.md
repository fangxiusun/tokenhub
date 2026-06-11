# rankings.go 代码阅读文档

## 1. 全局总结

该文件实现模型和供应商的排行榜功能，支持多种时间维度（今日/本周/本月/全部），提供模型排名、供应商排名、排名变动、历史趋势等数据。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `model` | 排行榜数据查询 |
| `sync` | 缓存锁 |

## 3. 类型定义

### `RankingsResponse`
排行榜响应：Models、Vendors、TopMovers、TopDroppers、历史数据

### `RankedModel`
模型排名：排名、名称、供应商、token 总量、份额、增长率

### `RankedVendor`
供应商排名：排名、名称、token 总量、份额、增长率、模型数量、Top 模型

### `RankingMover`
排名变动：模型、供应商、排名变化、当前排名、增长率

### `ModelHistorySeries` / `VendorShareSeries`
历史趋势数据

## 4. 函数详解

### `GetRankingsSnapshot(period) (*RankingsResponse, error)`
排行榜入口：
1. 检查缓存（5分钟 TTL）
2. 构建排行榜快照
3. 更新缓存

### `buildRankingsSnapshot(config, now)`
构建排行榜：
1. 查询当前时间段数据
2. 查询历史数据（计算增长率）
3. 构建模型排名
4. 构建供应商排名
5. 构建历史趋势
6. 计算排名变动

### 辅助函数
- `buildRankedModels` / `buildRankedVendors` — 排名构建
- `buildModelHistory` / `buildVendorShareHistory` — 历史趋势
- `buildRankingMovers` — 排名变动

## 5. 关键逻辑分析

1. **缓存策略**：5分钟缓存，避免频繁查询
2. **增长率计算**：`(current - previous) / previous * 100`
3. **Top N 限制**：模型排名 20，供应商 5，变动 6
4. **Others 归类**：超出 Top N 的归为 "Others"

## 6. 关联文件

- `model/ranking.go` — 排行榜数据查询
