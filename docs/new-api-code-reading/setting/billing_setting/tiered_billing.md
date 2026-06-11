# tiered_billing.go 代码阅读文档

## 1. 全局总结

该文件实现分层计费（Tiered Billing）配置管理，支持按模型设置不同的计费模式（固定倍率或表达式计费），并提供表达式的冒烟测试功能。

## 2. 依赖关系

- `fmt` — 错误格式化
- `github.com/QuantumNous/new-api/pkg/billingexpr` — 计费表达式引擎
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器
- `github.com/samber/lo` — 集合操作工具

## 3. 类型定义

| 常量名 | 值 | 说明 |
|--------|------|------|
| `BillingModeRatio` | `"ratio"` | 固定倍率计费模式 |
| `BillingModeTieredExpr` | `"tiered_expr"` | 表达式计费模式 |
| `BillingModeField` | `"billing_mode"` | 计费模式字段名 |
| `BillingExprField` | `"billing_expr"` | 计费表达式字段名 |

| 结构体 | 字段 | 说明 |
|--------|------|------|
| `BillingSetting` | `BillingMode map[string]string` | 模型→计费模式映射 |
| | `BillingExpr map[string]string` | 模型→计费表达式映射 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetBillingMode` | `func GetBillingMode(model string) string` | 获取模型的计费模式，默认返回 "ratio" |
| `GetBillingExpr` | `func GetBillingExpr(model string) (string, bool)` | 获取模型的计费表达式 |
| `GetBillingModeCopy` | `func GetBillingModeCopy() map[string]string` | 获取计费模式的深拷贝 |
| `GetBillingExprCopy` | `func GetBillingExprCopy() map[string]string` | 获取计费表达式的深拷贝 |
| `GetPricingSyncData` | `func GetPricingSyncData(base map[string]any) map[string]any` | 获取定价同步数据 |
| `SmokeTestExpr` | `func SmokeTestExpr(exprStr string) error` | 对计费表达式进行冒烟测试 |

## 5. 关键逻辑分析

- 通过 `init()` 注册到 `config.GlobalConfig`，键名为 `"billing_setting"`
- `SmokeTestExpr` 使用多组测试向量验证表达式：空输入、小规模、中等规模、大规模 token 参数
- 测试还包括不同的请求头和请求体组合
- `lo.Assign` 用于合并 map 并返回新 map，避免修改原 map

## 6. 关联文件

- `setting/config/config.go` — 全局配置管理器
- `pkg/billingexpr/` — 计费表达式引擎实现
- `controller/option.go` — 管理界面配置接口
