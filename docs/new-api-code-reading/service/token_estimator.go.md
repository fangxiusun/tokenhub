# token_estimator.go 代码阅读文档

## 1. 全局总结

该文件实现基于规则的 Token 估算器，支持 OpenAI、Gemini、Claude 三种厂商的估算权重。通过分析文本中的字符类型（字母、数字、CJK、Emoji 等）来估算 Token 数量。

## 2. 依赖关系

无外部依赖，纯算法实现。

## 3. 类型定义

### `Provider`
厂商类型：`OpenAI`、`Gemini`、`Claude`、`Unknown`

### `multipliers`
估算权重：
- `Word` — 英文单词
- `Number` — 数字
- `CJK` — 中日韩字符
- `Symbol` — 标点符号
- `MathSymbol` — 数学符号
- `URLDelim` — URL 分隔符
- `AtSign` — @符号
- `Emoji` — Emoji
- `Newline` / `Space` — 空白字符
- `BasePad` — 基础起步消耗

## 4. 函数详解

### `EstimateToken(provider, text) int`
Token 估算主函数：
1. 遍历每个 rune
2. 根据字符类型应用对应权重
3. 状态机跟踪单词边界
4. 向上取整 + BasePad

### `EstimateTokenByModel(model, text) int`
根据模型名称选择厂商

### 辅助函数
- `isCJK(r)` — CJK 字符判断
- `isEmoji(r)` — Emoji 判断
- `isMathSymbol(r)` — 数学符号判断
- `isURLDelim(r)` — URL 分隔符判断

## 5. 关键逻辑分析

1. **厂商差异**：不同厂商对同一字符的 token 消耗不同
2. **状态机**：跟踪字母/数字边界，避免重复计数
3. **数学符号**：Claude 的数学符号权重特别高（4.52）
4. **Emoji**：Claude 和 OpenAI 的 Emoji 权重较高

## 6. 关联文件

- `token_counter.go` — 使用估算器
- `tokenizer.go` — 精确 tokenizer
