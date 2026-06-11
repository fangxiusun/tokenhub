# user_usable_group.go 代码阅读文档

## 1. 全局总结

该文件管理用户可用分组配置，提供分组的增删查改功能，并支持并发安全访问。

## 2. 依赖关系

- `encoding/json` — JSON 序列化/反序列化
- `sync` — 读写锁
- `github.com/QuantumNous/new-api/common` — 系统日志

## 3. 类型定义

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `userUsableGroups` | `map[string]string` | 用户可用分组映射（分组名 → 描述） |
| `userUsableGroupsMutex` | `sync.RWMutex` | 读写锁 |

默认分组：
- `default` → "默认分组"
- `vip` → "vip分组"

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetUserUsableGroupsCopy` | `func GetUserUsableGroupsCopy() map[string]string` | 获取用户可用分组的深拷贝 |
| `UserUsableGroups2JSONString` | `func UserUsableGroups2JSONString() string` | 将分组配置序列化为 JSON |
| `UpdateUserUsableGroupsByJSONString` | `func UpdateUserUsableGroupsByJSONString(jsonStr string) error` | 从 JSON 更新分组配置 |
| `GetUsableGroupDescription` | `func GetUsableGroupDescription(groupName string) string` | 获取分组描述，不存在时返回分组名 |

## 5. 关键逻辑分析

- 使用 `sync.RWMutex` 保护并发访问
- `GetUserUsableGroupsCopy` 返回深拷贝，避免外部修改影响内部状态
- `UpdateUserUsableGroupsByJSONString` 会先清空再反序列化

## 6. 关联文件

- `model/user.go` — 用户模型，使用分组信息
- `setting/auto_group.go` — 自动分组配置
