# tool_billing.go 代码阅读文档

## 1. 全局总结

该文件实现工具调用的计费逻辑，包括 Web Search、File Search、Image Generation 等工具的费用计算。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 额度单位 |
| `operation_setting` | 工具价格配置 |

## 3. 类型定义

### `ToolCallUsage`
工具调用使用量：模型名称、Web Search/File Search 调用次数、Image Generation 信息

### `ToolCallItem`
单个工具计费项：名称、调用次数、单价、总价、额度

### `ToolCallResult`
工具调用计费结果：总额度、明细列表

## 4. 函数详解

### `ComputeToolCallQuota(usage, groupRatio) ToolCallResult`
计算工具调用总额度：
1. Web Search：`price * count / 1000 * groupRatio * quotaPerUnit`
2. File Search：同上
3. Image Generation：`price * groupRatio * quotaPerUnit`

## 5. 关键逻辑分析

1. **价格配置**：通过 `GetToolPriceForModel` 获取，支持模型前缀覆盖
2. **千次计价**：工具价格按千次调用计价
3. **Group 倍率**：应用分组倍率

## 6. 关联文件

- `text_quota.go` — 使用工具计费结果
- `operation_setting` — 工具价格配置
