# epay.go 代码阅读文档

## 1. 全局总结

该文件提供易支付（EPay）回调地址的获取功能。是一个非常简洁的工具文件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `operation_setting` | 自定义回调地址 |
| `system_setting` | 服务器地址 |

## 3. 函数详解

### `GetCallbackAddress() string`
- 优先返回自定义回调地址
- 未配置时返回服务器地址

## 4. 关联文件

- `setting/operation_setting` — 运营配置
- `setting/system_setting` — 系统配置
