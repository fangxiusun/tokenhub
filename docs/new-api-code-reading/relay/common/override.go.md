# override.go 代码阅读文档

## 1. 全局总结

本文件实现了 Relay 模块的参数覆盖（Param Override）系统，是整个 Relay 系统中最复杂的文件之一。支持对上游请求 JSON 的精细操作，包括 set/delete/move/copy/prepend/append/trim/replace/regex_replace 等 20+ 种操作模式，以及条件执行、通配符路径展开、请求头覆盖等高级功能。

## 2. 依赖关系

- `common`: JSON 操作
- `types`: 错误类型
- `gjson/sjson`: JSON 路径读写
- `lo`: 集合操作

## 3. 类型定义

### `ConditionOperation`
条件操作，用于条件执行参数覆盖：
- `Path`: JSON 路径
- `Mode`: 比较模式（full/prefix/suffix/contains/gt/gte/lt/lte）
- `Value`: 比较值
- `Invert`: 取反
- `PassMissingKey`: 键不存在时的行为

### `ParamOperation`
参数操作定义：
- `Path`: JSON 路径
- `Mode`: 操作模式（set/delete/move/copy/prepend/append/trim_prefix/trim_suffix/ensure_prefix/ensure_suffix/trim_space/to_lower/to_upper/replace/regex_replace/return_error/prune_objects/set_header/delete_header/copy_header/move_header/pass_headers/sync_fields）
- `Value`: 操作值
- `KeepOrigin`: 保留原值
- `From/To`: 源/目标路径
- `Conditions`: 条件列表
- `Logic`: 条件逻辑（AND/OR）

### `ParamOverrideReturnError`
参数覆盖返回错误类型，支持自定义 status_code、code、type、skip_retry。

## 4. 函数详解

### 核心入口

#### `ApplyParamOverride(jsonData, paramOverride, conditionContext) ([]byte, error)`
- 参数覆盖主入口，自动检测操作格式和旧格式

#### `ApplyParamOverrideWithRelayInfo(jsonData, info) ([]byte, error)`
- 基于 RelayInfo 的参数覆盖入口，自动构建上下文

### 操作执行

#### `applyOperations(jsonData, operations, conditionContext) ([]byte, error)`
- 在 []byte 上原地应用所有操作，避免 string 拷贝

#### `applyOperationsLegacy(jsonData, paramOverride, auditRecorder) ([]byte, error)`
- 旧格式兼容，按顶层 key 逐个写入

### 条件系统

#### `checkConditions(data, contextJSON, conditions, logic) (bool, error)`
- 检查条件列表，支持 AND/OR 逻辑

#### `checkSingleCondition(data, contextJSON, condition) (bool, error)`
- 检查单个条件，支持负数索引

### 路径展开

#### `expandWildcardPaths(data, path) ([]string, error)`
- 展开通配符路径（如 `contents.*.parts`）

### 头部覆盖

#### `setHeaderOverrideInContext(context, headerName, value, keepOrigin) error`
- 在上下文中设置请求头覆盖

#### `resolveHeaderOverrideValueByMapping(context, headerName, mapping) (string, bool, error)`
- 解析请求头覆盖的映射值（支持替换、追加、通配符）

## 5. 关键逻辑分析

1. **双格式兼容**: 同时支持旧格式（key-value 直接覆盖）和新格式（operations 数组）
2. **内存优化**: 全程在 []byte 上操作，避免 string 拷贝
3. **条件执行**: 每个操作可以有前置条件，支持 AND/OR 逻辑
4. **通配符路径**: 支持 `*` 通配符展开，批量操作多个路径
5. **请求头操作**: 支持 set_header/delete_header/copy_header/move_header/pass_headers
6. **sync_fields**: 在 JSON 路径和请求头之间同步字段值
7. **审计日志**: 敏感路径（model/messages/input 等）的操作会记录审计日志
8. **return_error**: 支持在参数覆盖阶段返回自定义错误，阻止请求发送

## 6. 关联文件

- `relay/channel/api_request.go`: Header Override 应用
- `relay/param_override_error.go`: 错误转换
