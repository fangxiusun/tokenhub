# types.go 代码阅读文档

## 1. 全局总结
该文件定义了性能指标系统的所有数据类型，包括存储接口、样本、查询参数、查询结果、桶键、计数器和原子桶等。提供类型安全的指标记录和查询。

## 2. 依赖关系
- **sync/atomic**: 原子操作。

## 3. 类型定义
### Store
接口，定义 Record 和 Query 方法。

### Sample
结构体，性能指标样本：模型、分组、延迟、TTFT、成功标志、输出 token 数、生成时间。

### QueryParams
结构体，查询参数：模型、分组、小时数。

### BucketPoint
结构体，桶点数据：时间戳、平均 TTFT、平均延迟、成功率、平均 TPS。

### GroupResult
结构体，分组结果：分组名、平均 TTFT、平均延迟、成功率、平均 TPS、时间序列。

### QueryResult
结构体，查询结果：模型名、序列模式、分组结果数组。

### ModelSummary
结构体，模型汇总：模型名、平均延迟、成功率、平均 TPS、请求数。

### SummaryAllResult
结构体，汇总结果：模型汇总数组。

### bucketKey
结构体，桶键：模型、分组、桶时间戳。

### counters
结构体，计数器：请求数、成功数、总延迟、TTFT 总和、TTFT 计数、输出 token 数、生成时间。

### atomicBucket
结构体，原子桶：使用 atomic.Int64 实现并发安全的计数器。

## 4. 函数详解
### atomicBucket.add
添加样本到原子桶：原子递增各计数器。

### atomicBucket.snapshot
获取桶的快照：原子加载各计数器值。

### atomicBucket.drain
排空桶：原子交换各计数器为 0，返回旧值。

### atomicBucket.addCounters
添加计数器到原子桶：原子递增各计数器。

## 5. 关键逻辑分析
- **原子操作**: 使用 atomic.Int64 确保并发安全的计数器操作。
- **接口设计**: Store 接口抽象存储实现，便于扩展。
- **JSON 标签**: 查询结果类型使用 JSON 标签，便于 API 响应序列化。

## 6. 关联文件
- **metrics.go**: 使用这些类型进行记录和查询。
- **flush.go**: 使用 counters 和 atomicBucket 类型。