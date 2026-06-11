# group.go 代码阅读文档

## 1. 全局总结

该文件提供用户分组管理功能，包括获取用户可用分组、自动分组、以及分组倍率查询。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `setting` | 分组配置 |
| `ratio_setting` | 倍率配置 |

## 3. 函数详解

### `GetUserUsableGroups(userGroup) map[string]string`
获取用户可用分组：
1. 获取默认可用分组
2. 应用特殊可用分组规则（`+:` 添加、`-:` 移除）
3. 确保用户自身分组在列表中

### `GroupInUserUsableGroups(userGroup, groupName) bool`
检查分组是否在用户可用分组中

### `GetUserAutoGroup(userGroup) []string`
获取用户可用的自动分组列表

### `GetUserGroupRatio(userGroup, group) float64`
获取用户使用某个分组的倍率：
1. 优先查找用户-分组特殊倍率
2. 回退到分组默认倍率

## 4. 关键逻辑分析

1. **特殊分组规则**：支持 `+:` 和 `-:` 前缀的添加/移除操作
2. **倍率优先级**：用户特殊倍率 > 分组默认倍率

## 5. 关联文件

- `setting` — 分组和倍率配置
- `channel_select.go` — 自动分组选择
