# return_path.go 代码阅读文档

## 1. 全局总结

该文件提供支付回调/返回路径的 URL 构建工具函数。

## 2. 依赖关系

- `common` — 主题感知路径
- `setting/system_setting` — 服务器地址

## 3. 类型定义

无。

## 4. 函数详解

### `paymentReturnPath(suffix string) string`
构建支付返回路径 URL。组合服务器地址和主题感知路径。

## 5. 关键逻辑分析

- 使用 `common.ThemeAwarePath` 支持不同前端主题的路径

## 6. 关联文件

- `common/` — 主题路径工具
