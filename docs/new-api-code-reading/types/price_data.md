# price_data.go 代码阅读文档

## 1. 全局概述

本文件定义了价格数据相关的结构体，用于存储和计算 AI 模型的计费信息。包括模型价格、各种比率（补全、缓存、图像、音频等）以及分组价格信息。

## 2. 依赖关系

- `fmt` — Go 标准库格式化包

## 3. 类型定义

### GroupRatioInfo 结构体

```go
type GroupRatioInfo struct {
    GroupRatio        float64
    GroupSpecialRatio float64
    HasSpecialRatio   bool
}
```

存储分组价格比率信息。

### PriceData 结构体

```go
type PriceData struct {
    FreeModel            bool
    ModelPrice           float64
    ModelRatio           float64
    CompletionRatio      float64
    CacheRatio           float64
    CacheCreationRatio   float64
    CacheCreation5mRatio float64
    CacheCreation1hRatio float64
    ImageRatio           float64
    AudioRatio           float64
    AudioCompletionRatio float64
    OtherRatios          map[string]float64
    UsePrice             bool
    Quota                int
    QuotaToPreConsume    int
    GroupRatioInfo       GroupRatioInfo
}
```

## 4. 函数详情

### AddOtherRatio

```go
func (p *PriceData) AddOtherRatio(key string, ratio float64)
```

添加额外的比率项。如果比率小于等于 0，则忽略。

### ToSetting

```go
func (p *PriceData) ToSetting() string
```

将 PriceData 转换为可读的设置字符串，用于日志和调试。

## 5. 关键逻辑分析

### 价格数据字段说明

| 字段 | 说明 |
|------|------|
| `FreeModel` | 是否为免费模型 |
| `ModelPrice` | 模型基础价格 |
| `ModelRatio` | 模型价格比率 |
| `CompletionRatio` | 补全价格比率（输出 Token 的价格倍率） |
| `CacheRatio` | 缓存价格比率 |
| `CacheCreationRatio` | 缓存创建价格比率 |
| `CacheCreation5mRatio` | 5 分钟缓存创建价格比率 |
| `CacheCreation1hRatio` | 1 小时缓存创建价格比率 |
| `ImageRatio` | 图像价格比率 |
| `AudioRatio` | 音频价格比率 |
| `AudioCompletionRatio` | 音频补全价格比率 |
| `OtherRatios` | 其他自定义比率 |
| `UsePrice` | 是否使用固定价格（而非比率计算） |
| `Quota` | 按次计费的最终额度（MJ / Task） |
| `QuotaToPreConsume` | 按量计费的预消耗额度 |
| `GroupRatioInfo` | 分组价格比率信息 |

### 计费逻辑

- **按量计费**：`ModelPrice × ModelRatio × GroupRatio × Token数量`
- **按次计费**：直接使用 `Quota` 字段的值
- **缓存优惠**：通过 `CacheRatio` 等字段实现缓存命中的价格折扣

### OtherRatios 扩展机制

`OtherRatios` 字段允许为特殊模型添加自定义比率，通过 `AddOtherRatio` 方法动态添加。

## 6. 相关文件

- `setting/model_setting.go` — 模型价格配置
- `setting/ratio.go` — 比率配置管理
- `relay/relay_info.go` — 中继信息中使用价格数据
- `model/log.go` — 日志中记录价格信息
