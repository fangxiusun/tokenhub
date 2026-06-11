# 内部工具包详细设计 (`pkg/`)

## 1. 概述
`pkg/` 目录包含为应用程序提供专业化功能的内部可重用包。这些包遵循内部包模式，具有清晰、专注的职责。

## 2. 子目录详细说明

### 2.1 `billingexpr/` - 计费表达式引擎
**职责**: 用于分层/动态定价的计费表达式引擎。

**核心功能**:
- 实现自定义表达式语言，用于根据 Token 使用量计算计费
- 支持变量：`p`（输入 Token）, `c`（输出 Token）, `b`（缓存 Token）
- 支持函数：`max()`, `min()`, `if()` 等
- 允许复杂的定价规则，如基于使用量的分层计费
- 表达式示例：`p*0.000002 + c*0.000008`（输入 $2/1M，输出 $8/1M）

**设计文档**:
- `pkg/billingexpr/expr.md` - 设计哲学、表达式语言、系统架构

### 2.2 `cachex/` - 混合缓存
**职责**: 结合 Redis 和内存缓存的混合缓存实现。

**核心功能**:
- 提供统一的缓存接口，带自动回退
- 支持基于 TTL 的过期和缓存失效
- 用于缓存模型比例、渠道配置和其他频繁访问的数据
- 写时失效，定时刷新策略

### 2.3 `ionet/` - io.net GPU 部署客户端
**职责**: io.net GPU 部署管理的客户端。

**核心功能**:
- 处理部署创建、管理和监控
- 提供硬件类型和位置查询
- 管理部署日志和容器信息
- 支持 GPU 资源调度和定价估算

### 2.4 `perf_metrics/` - 性能指标
**职责**: 性能指标收集和报告。

**核心功能**:
- 收集系统性能数据（CPU、内存、磁盘使用率）
- 支持不同时间桶（分钟、小时、天）的指标聚合
- 支持指标保留和清理
- 提供按需查询接口

## 3. 设计哲学
这些包遵循 **内部包模式**：
- 放置在 `pkg/` 中，表示它们是内部工具
- 具有清晰、专注的职责
- 设计为可在应用程序的不同部分重用
- 通过成为依赖树中的叶子包来避免循环依赖

## 4. 使用场景
- **billingexpr**: 在 `service/tiered_settle.go` 中用于分层计费结算
- **cachex**: 在 `model/channel_cache.go`, `model/user_cache.go` 中用于缓存管理
- **ionet**: 在 `controller/deployment.go` 中用于 GPU 部署管理
- **perf_metrics**: 在 `controller/perf_metrics.go` 中用于性能指标查询

---

## 关联文件列表

### 内部工具包核心文件
- `pkg/billingexpr/` - 计费表达式引擎目录
- `pkg/billingexpr/expr.md` - 设计文档
- `pkg/cachex/` - 混合缓存目录
- `pkg/ionet/` - io.net 客户端目录
- `pkg/perf_metrics/` - 性能指标目录

### 依赖此模块的文件
- `service/tiered_settle.go` - 分层结算（使用 billingexpr）
- `model/channel_cache.go` - 渠道缓存（使用 cachex）
- `model/user_cache.go` - 用户缓存（使用 cachex）
- `controller/deployment.go` - 部署管理（使用 ionet）
- `controller/perf_metrics.go` - 性能指标（使用 perf_metrics）

### 依赖的公共工具文件
- `common/redis.go` - Redis 工具
- `common/json.go` - JSON 序列化
- `common/constants.go` - 常量定义
