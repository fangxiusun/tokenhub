# price.go 代码阅读文档

## 1. 全局总结

本文件实现了 Relay 模块的价格计算逻辑，支持按 token 计费、按次计费和分层表达式计费三种模式。是计费系统的核心组件。

## 2. 依赖关系

- `model`: 用户模型
- `ratio_setting`: 模型价格、比率配置
- `billing_setting`: 计费模式、表达式
- `operation_setting`: 配额设置
- `pkg/billingexpr`: 分层计费表达式引擎

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `HandleGroupRatio(ctx, relayInfo) types.GroupRatioInfo`
- 处理分组比率，支持 auto_group 自动分组和用户组特殊比率

### `ModelPriceHelper(c, info, promptTokens, meta) (types.PriceData, error)`
- **按 token 计费**: preConsumedTokens × modelRatio × groupRatio
- **按价格计费**: modelPrice × QuotaPerUnit × groupRatio
- **分层计费**: 使用 billingexpr 引擎计算

### `ModelPriceHelperPerCall(c, info) (types.PriceData, error)`
- 按次计费（MJ、Task），使用固定价格或模型比率

### `HasModelBillingConfig(modelName) bool`
- 检查模型是否有计费配置

## 5. 关键逻辑分析

1. **三种计费模式**:
   - 按 token: promptTokens × ratio × groupRatio
   - 按价格: price × QuotaPerUnit × groupRatio
   - 分层表达式: billingexpr 引擎计算

2. **免费模型**: 当 groupRatio 或 modelPrice/ratio 为 0 时，跳过预扣费

3. **Claude 缓存价格**: 1h 缓存写入价格 = 5min 缓存价格 × (6/3.75)

4. **预扣费逻辑**: 取 max(promptTokens, PreConsumedQuota) + maxTokens 作为预扣基数

## 6. 关联文件

- `types/price_data.go`: PriceData 类型
- `pkg/billingexpr/`: 分层计费表达式引擎
- `setting/ratio_setting/`: 价格和比率配置
