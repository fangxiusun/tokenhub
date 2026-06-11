# config_test.go 代码阅读文档

## 1. 全局总结

该文件是 `config.go` 的单元测试，验证 `UpdateConfigFromMap` 函数在 map 替换、空 map 清除、标量字段不变等场景下的正确性。

## 2. 依赖关系

- `testing` — Go 标准测试框架

## 3. 类型定义

| 结构体 | 字段 | 说明 |
|--------|------|------|
| `testConfigWithMap` | `Modes map[string]string` | 测试用的 map 字段 |
| | `Exprs map[string]string` | 测试用的 map 字段 |
| | `Name string` | 测试用的标量字段 |

## 4. 函数详解

| 测试函数 | 说明 |
|----------|------|
| `TestUpdateConfigFromMap_MapReplacement` | 验证 map 字段的完整替换（移除旧键） |
| `TestUpdateConfigFromMap_EmptyMapClearsAll` | 验证空 map 能清空已有数据 |
| `TestUpdateConfigFromMap_ScalarFieldsUnchanged` | 验证未在 configMap 中的标量字段保持不变 |

## 5. 关键逻辑分析

- 测试重点验证 map 类型的"完全替换"语义：更新后的 map 不应包含旧数据
- `TestUpdateConfigFromMap_MapReplacement` 模拟删除 model-a 的场景
- `TestUpdateConfigFromMap_EmptyMapClearsAll` 验证 `{}` 能清空所有条目

## 6. 关联文件

- `setting/config/config.go` — 被测试的实现
