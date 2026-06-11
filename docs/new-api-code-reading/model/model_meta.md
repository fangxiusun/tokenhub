# model_meta.go 代码阅读文档

## 1. 全局总结
该文件定义了模型元数据（Model）的数据结构和相关的数据库操作函数。模型元数据用于描述 AI 模型的基本信息、状态、绑定渠道、计费类型等。该文件提供了模型的增删改查、搜索、统计等功能，是模型管理的核心数据层组件。

## 2. 依赖关系
- **strconv**: 用于字符串到整数的转换（在搜索功能中解析供应商 ID）。
- **strings**: 用于字符串处理（在 normalizeLookupValues 函数中进行修剪和去重）。
- **github.com/QuantumNous/new-api/common**: 提供通用工具函数（如 GetTimestamp、ChannelStatusEnabled）和跨数据库兼容性变量（如 commonGroupCol）。
- **gorm.io/gorm**: 用于数据库 ORM 操作，包括模型定义、查询构建、数据持久化等。
- **package model**: 当前包，提供数据库连接（DB）和其他模型定义。

## 3. 类型定义
### 常量定义
- **NameRuleExact (0)**: 精确匹配规则
- **NameRulePrefix (1)**: 前缀匹配规则
- **NameRuleContains (2)**: 包含匹配规则
- **NameRuleSuffix (3)**: 后缀匹配规则

### BoundChannel 结构体
表示模型绑定的渠道信息：
- `Name`: 渠道名称
- `Type`: 渠道类型

### Model 结构体
表示模型元数据的数据模型，包含以下字段：
- `Id`: 模型 ID（主键）
- `ModelName`: 模型名称（唯一索引，与 DeletedAt 联合唯一）
- `Description`: 模型描述
- `Icon`: 模型图标
- `Tags`: 模型标签
- `VendorID`: 供应商 ID（建立索引）
- `Endpoints`: 模型端点配置（JSON 格式）
- `Status`: 模型状态（默认 1）
- `SyncOfficial`: 是否同步官方模型（默认 1）
- `CreatedTime`: 创建时间戳
- `UpdatedTime`: 更新时间戳
- `DeletedAt`: 软删除时间戳（建立索引）
- `BoundChannels`: 绑定的渠道列表（非数据库字段）
- `EnableGroups`: 可用的用户组列表（非数据库字段）
- `QuotaTypes`: 计费类型列表（非数据库字段）
- `NameRule`: 名称匹配规则（默认 0）
- `MatchedModels`: 匹配的模型列表（非数据库字段）
- `MatchedCount`: 匹配数量（非数据库字段）

## 4. 函数详解
### 模型 CRUD 函数
- **(mi *Model) Insert() error**
  - 插入新的模型记录，设置创建时间和更新时间。
  - 特殊处理：由于 GORM 的 default 标签会覆盖零值，先创建记录，再更新 Status 和 SyncOfficial 字段以确保零值正确保存。

- **IsModelNameDuplicated(id int, name string) (bool, error)**
  - 检查模型名称是否已存在（排除指定 ID 的记录）。
  - 返回是否存在重复。

- **(mi *Model) Update() error**
  - 更新模型记录，强制更新所有字段（包括零值）。
  - 使用 Select 方法指定要更新的字段列表。

- **(mi *Model) Delete() error**
  - 软删除模型记录（GORM 的 DeletedAt 字段）。

### 查询函数
- **GetVendorModelCounts() (map[int64]int64, error)**
  - 统计每个供应商的模型数量。
  - 返回供应商 ID 到模型数量的映射。

- **GetAllModels(offset int, limit int) ([]*Model, error)**
  - 分页查询所有模型，按 ID 降序排列。
  - 返回模型列表。

- **GetBoundChannelsByModelsMap(modelNames []string) (map[string][]BoundChannel, error)**
  - 根据模型名称列表，查询每个模型绑定的渠道信息。
  - 通过 JOIN abilities 和 channels 表实现。
  - 返回模型名称到绑定渠道列表的映射。

- **GetPreferredModelOwnerChannelTypes(modelNames []string, groups []string) (map[string]int, error)**
  - 查询每个模型优先级最高的渠道类型。
  - 考虑渠道优先级、权重和启用状态。
  - 支持按用户组过滤。
  - 返回模型名称到渠道类型的映射。

- **SearchModels(keyword string, vendor string, offset int, limit int) ([]*Model, int64, error)**
  - 搜索模型：支持按关键词（模型名称、描述、标签）和供应商搜索。
  - 供应商参数可以是 ID 或名称。
  - 返回匹配的模型列表和总数。

### 辅助函数
- **normalizeLookupValues(values []string) []string**
  - 标准化查找值：修剪空格、去除空值、去重。
  - 用于清理输入参数。

## 5. 关键逻辑分析
- **软删除机制**: 使用 GORM 的 `DeletedAt` 字段实现软删除，通过唯一索引 `uk_model_name_delete_at` 确保删除后模型名称可重新使用。
- **零值处理**: 由于 GORM 的 default 标签会覆盖零值，`Insert` 和 `Update` 方法采用特殊策略确保 Status 和 SyncOfficial 等字段的零值能正确保存。
- **多表联查**: `GetBoundChannelsByModelsMap` 和 `GetPreferredModelOwnerChannelTypes` 通过 JOIN abilities 和 channels 表获取模型与渠道的关联信息。
- **优先级排序**: `GetPreferredModelOwnerChannelTypes` 按优先级、权重、渠道 ID 排序，确保返回优先级最高的渠道类型。
- **输入标准化**: `normalizeLookupValues` 函数清理输入参数，避免空值和重复值影响查询结果。

## 6. 关联文件
- **model/ability.go**: 可能包含 abilities 表的定义和操作（模型与渠道的绑定关系）。
- **model/channel.go**: 可能包含 channels 表的定义和操作（渠道信息）。
- **model/vendor.go**: 可能包含 vendors 表的定义和操作（供应商信息）。
- **common/constants.go**: 可能包含 ChannelStatusEnabled 等常量的定义。
- **common/main.go**: 可能包含 commonGroupCol 等跨数据库兼容性变量的定义。
- **controller/model.go**: 可能包含处理模型管理 HTTP 请求的控制器。
- **setting/model.go**: 可能包含模型配置的加载和更新逻辑。