# str.go 代码阅读文档

## 1. 全局总结

该文件提供字符串处理工具函数，包括 Sunday 搜索算法、去重、以及 Aho-Corasick 自动机的初始化和搜索。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `ahocorasick` | AC 自动机 |
| `sync` | 缓存锁 |

## 3. 函数详解

### `SundaySearch(text, pattern) bool`
Sunday 字符串搜索算法，用于单模式匹配

### `RemoveDuplicate(s []string) []string`
字符串切片去重

### `InitAc(dict []string) *goahocorasick.Machine`
初始化 AC 自动机

### `AcSearch(findText, dict, stopImmediately) (bool, []string)`
AC 自动机多模式搜索：
- 使用缓存的 AC 机器
- 支持立即停止模式
- 返回匹配的词列表

### `getOrBuildAC(dict) *goahocorasick.Machine`
获取或构建 AC 机器（带缓存）

## 4. 关键逻辑分析

1. **AC 缓存**：按字典内容哈希缓存 AC 机器
2. **大小写统一**：字典和搜索文本都转小写
3. **FNV 哈希**：使用 FNV-64a 对字典内容哈希

## 5. 关联文件

- `sensitive.go` — 敏感词检测
- `channel.go` — 关键词匹配
