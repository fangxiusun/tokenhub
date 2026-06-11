# return_path.go 代码阅读文档

## 1. 全局总结

该文件提供支付返回 URL 的生成功能。是一个非常简洁的工具文件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 主题感知路径 |
| `system_setting` | 服务器地址 |

## 3. 函数详解

### `PaymentReturnURL(suffix) string`
生成支付返回 URL：
- 拼接服务器地址 + 主题感知路径 + 后缀

## 4. 关联文件

- `setting/system_setting` — 服务器地址配置
