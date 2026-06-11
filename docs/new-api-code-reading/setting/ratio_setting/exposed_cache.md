# exposed_cache.go 代码阅读文档

## 1. 全局总结

该文件实现对外暴露数据的缓存机制，将模型倍率、价格等数据缓存 30 秒，减少重复计算。

## 2. 依赖关系

- `sync` — 互斥锁
- `sync/atomic` — 原子操作
- `time` — 时间处理
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `exposedCache` | `data` | `gin.H` | 缓存数据 |
| | `expiresAt` | `time.Time` | 过期时间 |

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `exposedData` | `atomic.Value` | 原子存储的缓存数据 |
| `rebuildMu` | `sync.Mutex` | 重建互斥锁 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `InvalidateExposedDataCache` | `func InvalidateExposedDataCache()` | 使缓存失效 |
| `cloneGinH` | `func cloneGinH(src gin.H) gin.H` | 深拷贝 gin.H |
| `GetExposedData` | `func GetExposedData() gin.H` | 获取暴露数据（带缓存） |

## 5. 关键逻辑分析

- 缓存 TTL 为 30 秒
- 使用双重检查锁定（Double-Checked Locking）模式
- `cloneGinH` 返回浅拷贝，避免缓存数据被外部修改
- 缓存数据包含：model_ratio、completion_ratio、cache_ratio、create_cache_ratio、model_price

## 6. 关联文件

- `setting/ratio_setting/model_ratio.go` — 提供倍率数据
- `setting/ratio_setting/cache_ratio.go` — 提供缓存倍率数据
- `controller/option.go` — 对外暴露接口
