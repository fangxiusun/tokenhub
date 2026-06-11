# pricing.go 代码阅读文档

## 1. 全局总结

`pricing.go` 是模型定价系统的核心文件，负责管理所有 AI 模型的定价信息。它从数据库中聚合模型元数据、供应商信息、渠道能力（Ability）和计费规则，构建全局定价缓存（`pricingMap`），供前端展示和后端计费使用。该文件实现了带缓存的定价查询机制，默认每分钟刷新一次，支持按前缀/后缀/包含规则匹配模型元数据，并维护模型支持的端点类型映射。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `encoding/json` | 解析模型自定义端点配置（JSON 格式） |
| `fmt` | 格式化日志输出 |
| `strings` | 字符串前缀/后缀/包含匹配，模型名规则判断 |
| `sync` | 互斥锁和读写锁，保护并发缓存访问 |
| `time` | 缓存过期时间判断（1分钟刷新间隔） |
| `github.com/QuantumNous/new-api/common` | 通用工具函数（端点类型判断、字符串包含检查等） |
| `github.com/QuantumNous/new-api/constant` | 常量定义（EndpointType 等） |
| `github.com/QuantumNous/new-api/setting/billing_setting` | 计费模式与计费表达式获取 |
| `github.com/QuantumNous/new-api/setting/ratio_setting` | 模型倍率、价格、缓存倍率等获取 |
| `github.com/QuantumNous/new-api/types` | 泛型 Set 数据结构 |

## 3. 类型定义

### Pricing 结构体
模型定价信息的核心结构，包含：
- `ModelName`：模型名称
- `Description`：模型描述
- `Icon`：模型图标
- `Tags`：标签（逗号分隔）
- `VendorID`：供应商 ID
- `QuotaType`：计费类型（0=按倍率，1=按固定价格）
- `ModelRatio`：输入倍率
- `ModelPrice`：固定价格
- `OwnerBy`：所有者
- `CompletionRatio`：输出倍率
- `CacheRatio`：缓存输入倍率（可选）
- `CreateCacheRatio`：创建缓存倍率（可选）
- `ImageRatio`：图像倍率（可选）
- `AudioRatio`/`AudioCompletionRatio`：音频倍率（可选）
- `EnableGroup`：启用的用户分组列表
- `SupportedEndpointTypes`：支持的端点类型列表
- `BillingMode`：计费模式（如 `tiered_expr`）
- `BillingExpr`：计费表达式
- `PricingVersion`：定价版本哈希

### PricingVendor 结构体
供应商信息，包含 `ID`、`Name`、`Description`、`Icon`，供前端展示。

### 全局变量
- `pricingMap`：全局定价列表缓存
- `vendorsList`：供应商列表缓存
- `supportedEndpointMap`：端点类型到路径的映射缓存
- `lastGetPricingTime`：上次刷新时间
- `updatePricingLock`：定价更新互斥锁
- `modelEnableGroups`：模型名到启用分组的快速查询缓存
- `modelQuotaTypeMap`：模型名到计费类型的快速查询缓存
- `modelSupportEndpointTypes`：模型到支持端点类型的映射缓存

## 4. 函数详解

### GetPricing() []Pricing
获取全局定价列表。采用双重检查锁模式：先检查缓存是否过期（>1分钟或为空），获取锁后再次检查，然后调用 `updatePricing()` 刷新。

### InvalidatePricingCache()
清除定价缓存，将 `pricingMap`、`vendorsList` 设为 nil，重置刷新时间。用于管理员手动刷新定价。

### GetVendors() []PricingVendor()
获取供应商列表。如果缓存为空会先触发 `GetPricing()` 刷新。

### GetModelSupportEndpointTypes(model string) []constant.EndpointType
查询指定模型支持的端点类型列表，使用读写锁保护并发访问。

### updatePricing()（私有）
核心刷新逻辑：
1. 获取所有启用的渠道能力（AbilityWithChannel）
2. 预加载模型元数据（Model），按精确/前缀/后缀/包含规则分类
3. 将非精确规则模型匹配到 `metaMap`
4. 预加载供应商，初始化默认供应商映射
5. 构建模型到分组的映射（使用 Set 去重）
6. 构建模型支持的端点类型列表（先填充原生端点，再用自定义端点覆盖）
7. 构建全局 `supportedEndpointMap`
8. 组装 `pricingMap`，包含倍率/价格/缓存倍率/计费表达式等
9. 刷新 `modelEnableGroups` 和 `modelQuotaTypeMap` 缓存

### GetSupportedEndpointMap() map[string]common.EndpointInfo
返回全局端点映射的只读引用。

## 5. 关键逻辑分析

**双重检查锁（Double-Check Locking）**：`GetPricing()` 使用该模式避免频繁加锁，仅在缓存过期时才获取互斥锁进行刷新。

**模型元数据匹配策略**：支持四种匹配规则——精确匹配、前缀匹配、后缀匹配、包含匹配。非精确规则按优先级依次匹配，已匹配的模型不会被重复覆盖。

**端点类型覆盖机制**：先根据渠道能力自动推导原生端点，再用模型表中自定义的 `Endpoints` JSON 配置覆盖，实现灵活的端点自定义。

**计费类型判断**：优先检查是否存在固定价格（`QuotaType=1`），否则使用倍率计费（`QuotaType=0`），支持缓存倍率、图像倍率、音频倍率等多种计费维度。

## 6. 关联文件

- `pricing_default.go`：默认供应商映射规则
- `pricing_refresh.go`：强制刷新定价的入口
- `model/model.go`：Model 结构体定义（NameRule 等）
- `model/vendor.go`：Vendor 结构体定义
- `model/ability.go`：AbilityWithChannel 结构体和 GetAllEnableAbilityWithChannels
- `setting/ratio_setting/`：模型倍率配置
- `setting/billing_setting/`：计费模式配置
- `common/endpoint.go`：端点信息和类型判断
