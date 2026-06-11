# gemini_generation_config_test.go 代码阅读文档

## 1. 全局摘要

该文件是 `GeminiChatGenerationConfig` 结构体的单元测试文件，验证自定义反序列化逻辑能正确处理 camelCase 和 snake_case 字段格式，并保留显式零值。

## 2. 依赖

- **测试框架**：
  - `testing`：Go 标准测试包
  - `github.com/stretchr/testify/assert`：断言库
  - `github.com/stretchr/testify/require`：必须断言库

- **项目内部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数

## 3. 类型定义

无独立类型定义，使用 `GeminiChatRequest` 和 `GeminiChatGenerationConfig` 结构体。

## 4. 函数详情

### TestGeminiChatGenerationConfigPreservesExplicitZeroValuesCamelCase()
```go
func TestGeminiChatGenerationConfigPreservesExplicitZeroValuesCamelCase(t *testing.T)
```
**功能**：测试 camelCase 字段格式的零值保留。

**测试场景**：
1. 解析包含 camelCase 字段（topP、topK、maxOutputTokens 等）的 JSON
2. 反序列化为 `GeminiChatRequest`
3. 重新序列化为 JSON
4. 验证字段存在且值为零值

### TestGeminiChatGenerationConfigPreservesExplicitZeroValuesSnakeCase()
```go
func TestGeminiChatGenerationConfigPreservesExplicitZeroValuesSnakeCase(t *testing.T)
```
**功能**：测试 snake_case 字段格式的零值保留。

**测试场景**：
1. 解析包含 snake_case 字段（top_p、top_k、max_output_tokens 等）的 JSON
2. 反序列化为 `GeminiChatRequest`
3. 重新序列化为 JSON
4. 验证字段存在且值为零值

## 5. 关键逻辑分析

1. **零值保留测试**：验证自定义反序列化逻辑能正确保留显式设置的零值，避免被 `omitempty` 过滤。

2. **双格式支持**：分别测试 camelCase 和 snake_case 两种 JSON 命名风格。

3. **字段完整性**：验证 topP、topK、maxOutputTokens、candidateCount、seed、responseLogprobs 等关键字段。

4. **往返测试**：验证 JSON → Go 结构体 → JSON 的序列化往返一致性。

## 6. 相关文件

- `dto/gemini.go`：被测试的 `GeminiChatRequest` 和 `GeminiChatGenerationConfig` 结构体
- `common/json.go`：JSON 工具函数