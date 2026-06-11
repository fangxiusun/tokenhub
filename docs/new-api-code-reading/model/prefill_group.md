# prefill_group.go 代码阅读文档

## 1. 全局总结
该文件定义了预填充组（PrefillGroup）的数据模型和相关的数据库操作函数。预填充组用于存储可复用的组信息，如模型组、标签组、端点组等，以便在前端下拉框中展示和选择。该文件提供了预填充组的增删改查功能，以及一个自定义的 JSON 值类型用于处理 JSON 字段。

## 2. 依赖关系
- **database/sql/driver**: 用于实现 `driver.Valuer` 接口（数据库写入）。
- **encoding/json**: 用于 JSON 序列化/反序列化。
- **github.com/QuantumNous/new-api/common**: 提供时间戳工具函数。
- **gorm.io/gorm**: 用于数据库 ORM 操作。

## 3. 类型定义
### JSONValue 类型
基于 `json.RawMessage` 实现的自定义 JSON 值类型，支持从数据库的 `[]byte` 和 `string` 两种类型读取：
- **Value() (driver.Value, error)**: 实现 `driver.Valuer` 接口，用于数据库写入。
- **Scan(value interface{}) error**: 实现 `sql.Scanner` 接口，兼容不同数据库驱动返回的类型。
- **MarshalJSON() ([]byte, error)**: 确保在对外编码时与 `json.RawMessage` 行为一致。
- **UnmarshalJSON(data []byte) error**: 确保在对外解码时与 `json.RawMessage` 行为一致。

### PrefillGroup 结构体
表示预填充组的数据模型，包含以下字段：
- `Id`: 组 ID（主键）
- `Name`: 组名称（唯一索引，排除已删除记录）
- `Type`: 组类型（如 model、tag、endpoint，建立索引）
- `Items`: 组项目列表（JSON 数组格式）
- `Description`: 组描述
- `CreatedTime`: 创建时间戳
- `UpdatedTime`: 更新时间戳
- `DeletedAt`: 软删除时间戳（建立索引）

## 4. 函数详解
### 插入函数
- **(g *PrefillGroup) Insert() error**
  - 插入新的预填充组记录。
  - 自动设置创建时间和更新时间。

### 查询函数
- **IsPrefillGroupNameDuplicated(id int, name string) (bool, error)**
  - 检查组名称是否已存在（排除指定 ID 的记录）。
  - 返回是否存在重复。

- **GetAllPrefillGroups(groupType string) ([]*PrefillGroup, error)**
  - 查询所有预填充组，可按类型过滤。
  - 按更新时间降序排列。

### 更新函数
- **(g *PrefillGroup) Update() error**
  - 更新预填充组记录。
  - 自动更新时间戳。

### 删除函数
- **DeletePrefillGroupByID(id int) error**
  - 根据 ID 软删除预填充组记录。

## 5. 关键逻辑分析
- **JSON 字段处理**: `JSONValue` 类型实现了 `driver.Valuer` 和 `sql.Scanner` 接口，确保 JSON 数据在数据库读写时的正确处理。同时实现了 `MarshalJSON` 和 `UnmarshalJSON` 方法，确保 JSON 编解码行为一致。
- **唯一索引**: 组名称的唯一索引使用条件 `where:deleted_at IS NULL`，确保只对未删除的记录进行唯一性检查，允许删除后重新使用相同名称。
- **软删除**: 使用 GORM 的 `DeletedAt` 字段实现软删除，通过唯一索引条件确保数据完整性。
- **时间戳管理**: 插入和更新操作自动设置时间戳，确保数据的可追溯性。

## 6. 关联文件
- **controller/prefill_group.go**: 可能包含处理预填充组 HTTP 请求的控制器。
- **router/prefill_group.go**: 可能包含预填充组相关的路由定义。
- **common/main.go**: 提供时间戳工具函数。
- **model/main.go**: 提供数据库连接和跨数据库兼容性工具。