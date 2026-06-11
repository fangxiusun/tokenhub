# status_code_ranges_test.go 代码阅读文档

## 1. 全局总结

该文件是 `status_code_ranges.go` 的单元测试，验证状态码范围解析、合并、禁用判断和重试判断的正确性。

## 2. 依赖关系

- `testing` — Go 标准测试框架
- `github.com/stretchr/testify/require` — 测试断言

## 3. 类型定义

无额外类型定义。

## 4. 函数详解

| 测试函数 | 说明 |
|----------|------|
| `TestParseHTTPStatusCodeRanges_CommaSeparated` | 验证逗号分隔的状态码解析 |
| `TestParseHTTPStatusCodeRanges_MergeAndNormalize` | 验证范围合并和排序 |
| `TestParseHTTPStatusCodeRanges_Invalid` | 验证无效输入的错误处理 |
| `TestParseHTTPStatusCodeRanges_NoComma_IsInvalid` | 验证空格分隔是无效的 |
| `TestShouldDisableByStatusCode` | 验证禁用判断逻辑 |
| `TestShouldRetryByStatusCode` | 验证重试判断逻辑 |
| `TestShouldRetryByStatusCode_DefaultMatchesLegacyBehavior` | 验证默认规则与旧版行为一致 |
| `TestIsAlwaysSkipRetryStatusCode` | 验证始终跳过重试的状态码 |

## 5. 关键逻辑分析

- 测试覆盖了正常的范围解析、边界情况和无效输入
- `TestShouldRetryByStatusCode_DefaultMatchesLegacyBehavior` 确保默认配置与历史行为兼容

## 6. 关联文件

- `setting/operation_setting/status_code_ranges.go` — 被测试的实现
