# openai_request_zero_value_test.go 代码阅读文档

## 1. 全局摘要

该文件是 OpenAI 请求结构体的单元测试文件，验证零值保留和系统角色名称获取逻辑的正确性。

## 2. 依赖

- **测试框架**：
  - `testing`：Go 标准测试包
  - `github.com/stretchr/testify/require`：必须断言库
  - `github.com/tidwall/gjson`：JSON 查询库

- **项目内部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数

## 3. 类型定义

无独立类型定义。

## 4. 函数详情

### TestGeneralOpenAIRequestPreserveExplicitZeroValues()
```go
func TestGeneralOpenAIRequestPreserveExplicitZeroValues(t *testing.T)
```
**功能**：测试 `GeneralOpenAIRequest` 的零值保留。

**测试场景**：
1. 解析包含显式零值的 JSON（stream=false、max_tokens=0 等）
2. 反序列化为 `GeneralOpenAIRequest`
3. 重新序列化为 JSON
4. 验证所有字段在输出 JSON 中存在

**验证字段**：stream、max_tokens、max_completion_tokens、top_p、top_k、n、frequency_penalty、presence_penalty、seed、logprobs、top_logprobs、dimensions、return_images、return_related_questions

### TestOpenAIResponsesRequestPreserveExplicitZeroValues()
```go
func TestOpenAIResponsesRequestPreserveExplicitZeroValues(t *testing.T)
```
**功能**：测试 `OpenAIResponsesRequest` 的零值保留。

**测试场景**：
1. 解析包含显式零值的 JSON
2. 验证字段在输出 JSON 中存在

**验证字段**：max_output_tokens、max_tool_calls、stream、top_p

### TestGeneralOpenAIRequestGetSystemRoleName()
```go
func TestGeneralOpenAIRequestGetSystemRoleName(t *testing.T)
```
**功能**：测试系统角色名称获取逻辑。

**测试用例**：
1. `o1` → "developer"
2. `o3-mini-high` → "developer"
3. `o4-mini` → "developer"
4. `o1-mini` → "system"
5. `o1-preview` → "system"
6. `gpt-5` → "developer"
7. `omni-moderation-latest` → "system"

## 5. 关键逻辑分析

1. **零值保留**：验证使用指针类型和 `omitempty` 标签的字段能正确保留显式设置的零值。

2. **系统角色适配**：验证不同模型的系统角色名称适配逻辑：
   - o1/o3/o4 系列（除 mini/preview）→ "developer"
   - o1-mini/o1-preview → "system"
   - gpt-5 系列 → "developer"
   - 其他模型 → "system"

3. **字段存在性检查**：使用 `gjson.GetBytes()` 检查字段在序列化后的 JSON 中是否存在。

4. **表驱动测试**：使用表驱动测试模式，覆盖多种模型场景。

## 6. 相关文件

- `dto/openai_request.go`：被测试的 `GeneralOpenAIRequest` 和 `OpenAIResponsesRequest` 结构体
- `common/json.go`：JSON 工具函数