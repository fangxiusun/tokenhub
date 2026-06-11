# group.go 代码阅读文档

## 1. 全局总结

该文件实现了用户组（Group）的查询接口，包括获取所有组名列表和获取当前用户可用的组信息（含倍率和描述）。

## 2. 依赖关系

- `model` — 用户组模型
- `service` — 用户可用组查询、倍率计算
- `setting` — 组描述获取
- `setting/ratio_setting` — 组倍率配置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetGroups(c *gin.Context)`
获取所有组名列表。从 `ratio_setting.GetGroupRatioCopy()` 中提取组名。

### `GetUserGroups(c *gin.Context)`
获取当前用户可用的组信息。返回每个可用组的倍率和描述，支持 auto 组（自动分配）。

## 5. 关键逻辑分析

- 组信息从 `ratio_setting` 配置中读取
- auto 组的倍率显示为 "自动"
- 只返回用户实际可用的组（通过 `service.GetUserUsableGroups` 过滤）

## 6. 关联文件

- `setting/ratio_setting/` — 组倍率配置
- `service/user.go` — 用户可用组查询
