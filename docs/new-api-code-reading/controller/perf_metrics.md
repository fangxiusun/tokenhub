# perf_metrics.go 代码阅读文档

## 1. 全局总结

该文件提供了性能指标查询接口，包括汇总指标和按模型/组的详细指标。

## 2. 依赖关系

- `pkg/perf_metrics` — 性能指标查询
- `setting/ratio_setting` — 组配置
- `gin-gonic/gin` — HTTP 框架
- `samber/lo` — 集合操作

## 3. 类型定义

无。

## 4. 函数详解

### `GetPerfMetricsSummary(c *gin.Context)`
获取所有活跃组的性能指标汇总。默认查询最近 24 小时。

### `GetPerfMetrics(c *gin.Context)`
获取指定模型的性能指标。支持按组和时间范围过滤。

### `filterActiveGroups(groups)`
过滤掉不存在于活跃组配置中的结果。

## 5. 关键逻辑分析

- 活跃组 = 配置中的组 + "auto"
- 时间范围通过 `hours` 参数控制，默认 24 小时

## 6. 关联文件

- `pkg/perf_metrics/` — 性能指标实现
