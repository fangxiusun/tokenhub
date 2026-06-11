# violation_fee.go 代码阅读文档

## 1. 全局总结

该文件实现违规费用扣除功能，当用户请求触发安全检查（如 CSAM 违规）时，额外扣除费用作为惩罚。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 额度单位 |
| `logger` | 日志 |
| `model` | 额度操作 |
| `relaycommon` | RelayInfo |
| `model_setting` | Grok 设置 |
| `types` | 错误类型 |
| `decimal` | 精确计算 |

## 3. 类型定义

常量：
- `ViolationFeeCodePrefix = "violation_fee."` — 违规费用代码前缀
- `CSAMViolationMarker` / `ContentViolatesUsageMarker` — 违规标记

## 4. 函数详解

### `IsViolationFeeCode(code) bool`
检查是否为违规费用代码

### `HasCSAMViolationMarker(err) bool`
检查错误是否包含 CSAM 违规标记

### `NormalizeViolationFeeError(err) *types.NewAPIError`
标准化违规费用错误：
- CSAM 标记 → 设置稳定错误代码 + 跳过重试
- 已有违规前缀 → 启用跳过重试

### `ChargeViolationFeeIfNeeded(ctx, relayInfo, apiErr) bool`
违规费用扣除主逻辑：
1. 检查是否应该收取违规费用
2. 获取 Grok 设置
3. 计算费用额度
4. 扣除费用
5. 记录日志

### `calcViolationFeeQuota(amount, groupRatio) int`
计算违规费用额度：`amount * quotaPerUnit * groupRatio`

## 5. 关键逻辑分析

1. **Grok 特殊处理**：使用 Grok 的费用策略
2. **跳过重试**：违规费用错误不重试
3. **费用计算**：基于金额和分组倍率

## 6. 关联文件

- `model_setting/grok.go` — Grok 违规设置
- `quota.go` — PostConsumeQuota
