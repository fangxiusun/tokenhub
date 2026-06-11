# expose_ratio.go 代码阅读文档

## 1. 全局总结

该文件控制倍率数据是否对外暴露（通过 API 返回给前端）。

## 2. 依赖关系

- `sync/atomic` — 原子操作

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `exposeRatioEnabled` | `atomic.Bool` | `false` | 是否启用倍率暴露 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `SetExposeRatioEnabled` | `func SetExposeRatioEnabled(enabled bool)` | 设置倍率暴露开关 |
| `IsExposeRatioEnabled` | `func IsExposeRatioEnabled() bool` | 检查倍率是否对外暴露 |

## 5. 关键逻辑分析

- 使用 `atomic.Bool` 实现无锁读写
- 默认关闭，需管理员手动开启

## 6. 关联文件

- `setting/ratio_setting/exposed_cache.go` — 暴露数据缓存
- `controller/option.go` — 管理界面配置
