# pricing_default.go 代码阅读文档

## 1. 全局总结

`pricing_default.go` 负责模型定价系统中的默认供应商自动识别与创建。当数据库中没有为某个模型手动配置供应商元数据时，该文件通过预定义的前缀匹配规则自动识别模型所属供应商，并创建对应的供应商记录和模型元数据。它还维护了供应商的默认图标映射。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `strings` | 字符串小写转换和包含匹配 |

该文件无外部项目依赖，是纯逻辑文件。

## 3. 类型定义

### defaultVendorRules（变量）
`map[string]string` 类型，定义模型名前缀到供应商名称的映射规则。包含 38 条规则，覆盖 OpenAI、Anthropic、Google、Moonshot、智谱、阿里巴巴、DeepSeek、MiniMax、百度、讯飞、腾讯、Cohere、Cloudflare、360、零一万物、Jina、Mistral、xAI、Meta、字节跳动、快手、即梦、Vidu 等供应商。

### defaultVendorIcons（变量）
`map[string]string` 类型，定义供应商名称到图标的映射。供新建供应商时自动设置图标。

## 4. 函数详解

### initDefaultVendorMapping(metaMap map[string]*Model, vendorMap map[int]*Vendor, enableAbilities []AbilityWithChannel)
核心函数，由 `updatePricing()` 调用。遍历所有已启用的渠道能力，对每个模型：
1. 如果 `metaMap` 中已存在该模型（已有手动配置），跳过
2. 否则将模型名转为小写，遍历 `defaultVendorRules` 进行前缀匹配
3. 匹配成功后调用 `getOrCreateVendor` 查找或创建供应商
4. 创建默认模型元数据（`NameRule` 为精确匹配，`Status=1`）

### getOrCreateVendor(vendorName string, vendorMap map[int]*Vendor) int
查找或创建供应商：
1. 先在 `vendorMap` 中查找是否已存在同名供应商
2. 若不存在，创建新供应商并插入数据库
3. 设置默认图标
4. 返回供应商 ID

### getDefaultVendorIcon(vendorName string) string
从 `defaultVendorIcons` 映射中获取供应商的默认图标名称，不存在则返回空字符串。

## 5. 关键逻辑分析

**前缀匹配策略**：匹配规则按 Map 遍历顺序执行，第一个命中即停止。规则设计时需注意重叠问题（如 `glm-` 和 `chatglm` 同属智谱）。

**幂等性设计**：`initDefaultVendorMapping` 只处理 `metaMap` 中不存在的模型，已手动配置的模型不会被覆盖，保证管理员的手动设置优先。

**供应商自动创建**：新发现的模型会自动创建供应商记录，降低了系统初始化和接入新模型的运维成本。

## 6. 关联文件

- `pricing.go`：在 `updatePricing()` 中调用 `initDefaultVendorMapping`
- `model/vendor.go`：`Vendor` 结构体定义和 `Insert()` 方法
- `model/model.go`：`Model` 结构体定义（`NameRule` 等）
