# auto_group.go 代码阅读文档

## 1. 全局总结

该文件管理自动分组配置，提供自动分组的增删查改功能。自动分组用于用户注册时自动分配的用户组。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/common` — JSON 序列化/反序列化

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `autoGroups` | `[]string` | 自动分组列表，默认包含 "default" |
| `DefaultUseAutoGroup` | `bool` | 是否默认使用自动分组，默认 `false` |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `ContainsAutoGroup` | `func ContainsAutoGroup(group string) bool` | 检查指定分组是否在自动分组列表中 |
| `UpdateAutoGroupsByJsonString` | `func UpdateAutoGroupsByJsonString(jsonString string) error` | 从 JSON 字符串更新自动分组列表 |
| `AutoGroups2JsonString` | `func AutoGroups2JsonString() string` | 将自动分组列表序列化为 JSON 字符串 |
| `GetAutoGroups` | `func GetAutoGroups() []string` | 获取自动分组列表 |

## 5. 关键逻辑分析

- `autoGroups` 是包级变量，直接存储当前配置
- `UpdateAutoGroupsByJsonString` 会先清空再反序列化，确保完全替换
- `ContainsAutoGroup` 使用线性搜索，适用于小规模列表

## 6. 关联文件

- `model/user.go` — 用户注册时使用自动分组
- `common/json.go` — JSON 操作封装
