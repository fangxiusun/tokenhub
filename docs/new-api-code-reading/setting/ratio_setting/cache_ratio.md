# cache_ratio.go 代码阅读文档

## 1. 全局总结

该文件管理缓存倍率（Cache Ratio）配置，包括缓存读取倍率和缓存创建倍率，用于计算缓存命中时的计费折扣。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/types` — RWMap 并发安全 Map

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `defaultCacheRatio` | `map[string]float64` | 默认缓存读取倍率（如 GPT-4o 为 0.5） |
| `defaultCreateCacheRatio` | `map[string]float64` | 默认缓存创建倍率（Claude 系列为 1.25） |
| `cacheRatioMap` | `types.RWMap[string, float64]` | 缓存读取倍率映射 |
| `createCacheRatioMap` | `types.RWMap[string, float64]` | 缓存创建倍率映射 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetCacheRatioMap` | `func GetCacheRatioMap() map[string]float64` | 获取缓存读取倍率映射 |
| `CacheRatio2JSONString` | `func CacheRatio2JSONString() string` | 序列化缓存读取倍率 |
| `CreateCacheRatio2JSONString` | `func CreateCacheRatio2JSONString() string` | 序列化缓存创建倍率 |
| `UpdateCacheRatioByJSONString` | `func UpdateCacheRatioByJSONString(jsonStr string) error` | 从 JSON 更新缓存读取倍率 |
| `UpdateCreateCacheRatioByJSONString` | `func UpdateCreateCacheRatioByJSONString(jsonStr string) error` | 从 JSON 更新缓存创建倍率 |
| `GetCacheRatio` | `func GetCacheRatio(name string) (float64, bool)` | 获取模型的缓存读取倍率 |
| `GetCreateCacheRatio` | `func GetCreateCacheRatio(name string) (float64, bool)` | 获取模型的缓存创建倍率 |
| `GetCacheRatioCopy` | `func GetCacheRatioCopy() map[string]float64` | 获取缓存读取倍率的深拷贝 |
| `GetCreateCacheRatioCopy` | `func GetCreateCacheRatioCopy() map[string]float64` | 获取缓存创建倍率的深拷贝 |

## 5. 关键逻辑分析

- 缓存读取倍率默认 1（未配置时），表示无折扣
- 缓存创建倍率默认 1.25，表示创建缓存时额外收费 25%
- Claude 系列模型统一使用 1.25 的创建倍率
- 更新操作会触发 `InvalidateExposedDataCache` 刷新暴露数据缓存

## 6. 关联文件

- `setting/ratio_setting/model_ratio.go` — 模型倍率配置
- `setting/ratio_setting/exposed_cache.go` — 暴露数据缓存
- `service/billing.go` — 计费服务
