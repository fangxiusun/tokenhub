# midjourney.go 代码阅读文档

## 1. 全局总结

该文件定义 Midjourney 相关的功能开关配置。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `MjNotifyEnabled` | `bool` | `false` | 是否启用 Midjourney 通知 |
| `MjAccountFilterEnabled` | `bool` | `false` | 是否启用账号过滤 |
| `MjModeClearEnabled` | `bool` | `false` | 是否启用模式清理 |
| `MjForwardUrlEnabled` | `bool` | `true` | 是否启用转发 URL |
| `MjActionCheckSuccessEnabled` | `bool` | `true` | 是否启用操作成功检查 |

## 4. 函数详解

无函数定义，仅变量声明。

## 5. 关键逻辑分析

所有变量均为包级全局变量，通过其他模块直接读取。

## 6. 关联文件

- `controller/midjourney.go` — Midjourney 相关控制器
- `relay/midjourney/` — Midjourney 中继实现
