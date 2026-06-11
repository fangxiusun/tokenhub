# group_ratio.go 代码阅读文档

## 1. 全局总结

该文件管理用户组倍率配置，支持按用户组设置不同的倍率，以及跨组倍率和特殊可用组配置。

## 2. 依赖关系

- `encoding/json` — JSON 序列化
- `errors` — 错误创建
- `github.com/QuantumNous/new-api/common` — 系统日志
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器
- `github.com/QuantumNous/new-api/types` — RWMap 并发安全 Map

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `GroupRatioSetting` | `GroupRatio` | `*types.RWMap[string, float64]` | 用户组倍率 |
| | `GroupGroupRatio` | `*types.RWMap[string, map[string]float64]` | 跨组倍率 |
| | `GroupSpecialUsableGroup` | `*types.RWMap[string, map[string]string]` | 特殊可用组 |

默认分组倍率：`default=1`、`vip=1`、`svip=1`。

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetGroupRatioSetting` | `func GetGroupRatioSetting() *GroupRatioSetting` | 获取组倍率设置 |
| `GetGroupRatioCopy` | `func GetGroupRatioCopy() map[string]float64` | 获取组倍率的深拷贝 |
| `ContainsGroupRatio` | `func ContainsGroupRatio(name string) bool` | 检查是否包含指定组倍率 |
| `GroupRatio2JSONString` | `func GroupRatio2JSONString() string` | 序列化组倍率 |
| `UpdateGroupRatioByJSONString` | `func UpdateGroupRatioByJSONString(jsonStr string) error` | 从 JSON 更新组倍率 |
| `GetGroupRatio` | `func GetGroupRatio(name string) float64` | 获取指定组的倍率 |
| `GetGroupGroupRatio` | `func GetGroupGroupRatio(userGroup, usingGroup string) (float64, bool)` | 获取跨组倍率 |
| `CheckGroupRatio` | `func CheckGroupRatio(jsonStr string) error` | 校验组倍率配置 |

## 5. 关键逻辑分析

- 组倍率默认为 1，未找到时也返回 1
- 跨组倍率用于用户组使用其他组的资源时的倍率调整
- 特殊可用组支持追加（`append_1`）和移除（`-:remove_1`）语法
- `CheckGroupRatio` 校验倍率不能为负数

## 6. 关联文件

- `setting/ratio_setting/model_ratio.go` — 模型倍率
- `service/billing.go` — 计费服务
- `model/user.go` — 用户模型
