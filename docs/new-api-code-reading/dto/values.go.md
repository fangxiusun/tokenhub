# values.go 代码阅读文档

## 1. 全局摘要

该文件定义了自定义值类型 `StringValue`、`IntValue`、`BoolValue`，支持从字符串或其他类型进行 JSON 反序列化。用于处理上游 API 返回的类型不一致的数据。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `strconv`：字符串转换

## 3. 类型定义

### StringValue 类型
自定义字符串类型，支持从字符串或数字反序列化：
```go
type StringValue string
```

### IntValue 类型
自定义整数类型，支持从整数或字符串反序列化：
```go
type IntValue int
```

### BoolValue 类型
自定义布尔类型，支持从布尔值或字符串反序列化：
```go
type BoolValue bool
```

## 4. 函数详情

### StringValue 方法

**UnmarshalJSON(data []byte) error**：自定义反序列化：
1. 尝试解析为字符串
2. 尝试解析为数字（转换为字符串）
3. 回退到标准字符串解析

**MarshalJSON() ([]byte, error)**：自定义序列化。

### IntValue 方法

**UnmarshalJSON(b []byte) error**：自定义反序列化：
1. 尝试解析为整数
2. 尝试解析为字符串（转换为整数）

**MarshalJSON() ([]byte, error)**：自定义序列化。

### BoolValue 方法

**UnmarshalJSON(data []byte) error**：自定义反序列化：
1. 尝试解析为布尔值
2. 尝试解析为字符串（"true"/"false"）

**MarshalJSON() ([]byte, error)**：自定义序列化。

## 5. 关键逻辑分析

1. **类型兼容性**：处理上游 API 返回类型不一致的情况（如数字作为字符串返回）。

2. **优雅降级**：反序列化失败时尝试其他类型，提高容错性。

3. **序列化一致性**：序列化时使用标准类型，确保输出格式统一。

## 6. 相关文件

- `common/json.go`：JSON 工具函数
- `relay/`：中继适配器使用这些类型处理上游响应