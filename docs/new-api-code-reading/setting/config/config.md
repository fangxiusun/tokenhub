# config.go 代码阅读文档

## 1. 全局总结

该文件实现统一的配置管理器（ConfigManager），支持配置模块的注册、数据库加载/保存、以及基于反射的结构体与 map 之间的转换。是整个 setting 包的核心基础设施。

## 2. 依赖关系

- `encoding/json` — JSON 序列化/反序列化
- `reflect` — 反射操作
- `strconv` — 类型转换
- `strings` — 字符串处理
- `sync` — 互斥锁
- `github.com/QuantumNous/new-api/common` — 系统日志

## 3. 类型定义

| 结构体 | 字段 | 说明 |
|--------|------|------|
| `ConfigManager` | `configs map[string]interface{}` | 已注册的配置模块映射 |
| | `mutex sync.RWMutex` | 并发读写锁 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `NewConfigManager` | `func NewConfigManager() *ConfigManager` | 创建新的配置管理器 |
| `Register` | `func (cm *ConfigManager) Register(name string, config interface{})` | 注册配置模块 |
| `Get` | `func (cm *ConfigManager) Get(name string) interface{}` | 获取配置模块 |
| `LoadFromDB` | `func (cm *ConfigManager) LoadFromDB(options map[string]string) error` | 从数据库加载配置 |
| `SaveToDB` | `func (cm *ConfigManager) SaveToDB(updateFunc func(key, value string) error) error` | 保存配置到数据库 |
| `ExportAllConfigs` | `func (cm *ConfigManager) ExportAllConfigs() map[string]string` | 导出所有配置为扁平结构 |
| `configToMap` | `func configToMap(config interface{}) (map[string]string, error)` | 结构体转 map |
| `updateConfigFromMap` | `func updateConfigFromMap(config interface{}, configMap map[string]string) error` | 从 map 更新结构体 |
| `ConfigToMap` | `func ConfigToMap(config interface{}) (map[string]string, error)` | 导出的 configToMap |
| `UpdateConfigFromMap` | `func UpdateConfigFromMap(config interface{}, configMap map[string]string) error` | 导出的 updateConfigFromMap |

## 5. 关键逻辑分析

- `GlobalConfig` 是全局单例配置管理器
- `Register` 在 `init()` 中调用，将配置模块注册到管理器
- `LoadFromDB` 使用 `"模块名.配置项"` 格式的键名前缀匹配
- `configToMap` 使用反射遍历结构体字段，支持 string/bool/int/uint/float/ptr/map/slice/struct 类型
- `updateConfigFromMap` 中 map 类型会先分配新 map 再反序列化，确保旧键被正确删除
- `updateConfigFromMap` 对 int 类型兼容 float 格式字符串（如 "2.000000"）

## 6. 关联文件

- `setting/billing_setting/tiered_billing.go` — 注册示例
- `setting/model_setting/global.go` — 注册示例
- `model/option.go` — 数据库存储
