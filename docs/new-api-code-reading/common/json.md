# json.go 代码阅读文档

## 1. 全局总结

`json.go` 是 JSON 操作的封装文件，提供了对标准库 `encoding/json` 的统一包装。根据项目规范（AGENTS.md Rule 1），所有业务代码中的 JSON 序列化/反序列化操作都必须使用此文件提供的函数，禁止直接导入和使用 `encoding/json`。这种设计便于未来替换为更高性能的 JSON 库。

## 2. 依赖关系

### 标准库依赖
- `bytes` - 字节切片处理
- `encoding/json` - JSON 编解码
- `io` - I/O 流操作

### 项目内部依赖
- 无（但被项目中几乎所有文件依赖）

## 3. 类型定义

该文件未定义新的类型，但使用并暴露了 `encoding/json` 中的类型（如 `json.RawMessage`）。

## 4. 函数详解

### Unmarshal(data []byte, v any) error
```go
func Unmarshal(data []byte, v any) error
```
将字节切片反序列化到指定变量。

**参数：**
- `data` - JSON 字节数据
- `v` - 目标变量指针

**返回值：**
- `error` - 反序列化错误

### UnmarshalJsonStr(data string, v any) error
```go
func UnmarshalJsonStr(data string, v any) error
```
将字符串形式的 JSON 反序列化到指定变量。

**参数：**
- `data` - JSON 字符串
- `v` - 目标变量指针

**返回值：**
- `error` - 反序列化错误

**实现逻辑：**
- 使用 `StringToByteSlice()` 辅助函数将字符串转换为字节切片
- 避免显式的字符串到字节的转换开销

### DecodeJson(reader io.Reader, v any) error
```go
func DecodeJson(reader io.Reader, v any) error
```
从 I/O Reader 流中解码 JSON。

**参数：**
- `reader` - 数据源
- `v` - 目标变量指针

**返回值：**
- `error` - 解码错误

**使用场景：**
- HTTP 请求体解析
- 文件流读取

### Marshal(v any) ([]byte, error)
```go
func Marshal(v any) ([]byte, error)
```
将变量序列化为 JSON 字节切片。

**参数：**
- `v` - 待序列化的变量

**返回值：**
- `[]byte` - JSON 字节数据
- `error` - 序列化错误

### GetJsonType(data json.RawMessage) string
```go
func GetJsonType(data json.RawMessage) string
```
获取 JSON 数据的类型标识。

**参数：**
- `data` - JSON 原始数据

**返回值：**
- `string` - 类型标识

**返回值映射：**
| 首字符 | 返回类型 |
|--------|---------|
| `{` | "object" |
| `[` | "array" |
| `"` | "string" |
| `t`, `f` | "boolean" |
| `n` | "null" |
| 其他 | "number" |
| 空数据 | "unknown" |

### JsonRawMessageToString(data json.RawMessage) string
```go
func JsonRawMessageToString(data json.RawMessage) string
```
将 JSON 原始消息转换为字符串表示。

**参数：**
- `data` - JSON 原始数据

**返回值：**
- `string` - 格式化后的字符串

**实现逻辑：**
1. 去除空白字符
2. 处理空数据和 null 值，返回空字符串
3. 如果不是字符串类型，直接返回原始文本
4. 如果是字符串类型，进行反序列化获取实际值

## 5. 关键逻辑分析

### 统一接口设计
```go
// 项目规范要求：所有 JSON 操作必须使用这些函数
common.Marshal(v)
common.Unmarshal(data, &v)
common.UnmarshalJsonStr(str, &v)
common.DecodeJson(reader, &v)
```

### 类型检测实现
```go
switch firstChar {
case '{':
    return "object"
case '[':
    return "array"
// ...
}
```
- 通过检查 JSON 数据的第一个非空白字符判断类型
- 高效的 O(1) 时间复杂度

### 字符串解码处理
```go
if trimmed[0] != '"' {
    return string(trimmed)
}
var value string
if err := Unmarshal(trimmed, &value); err != nil {
    return string(trimmed)
}
return value
```
- 只对字符串类型进行特殊处理
- 解码失败时降级为原始文本

## 6. 关联文件

- `common/json_test.go` - 单元测试文件
- `common/str_utils.go` - 包含 `StringToByteSlice()` 辅助函数
- `relay/` - 所有中继模块使用这些 JSON 函数进行数据转换
- `controller/` - 所有控制器使用这些函数处理请求/响应
- `model/` - 数据模型序列化/反序列化
