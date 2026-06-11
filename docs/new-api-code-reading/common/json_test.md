# json_test.go 代码阅读文档

## 1. 全局总结

`json_test.go` 是 `common` 包中 JSON 相关函数的单元测试文件。目前只包含 `JsonRawMessageToString` 函数的测试用例，覆盖了 JSON 对象、字符串、null 和空值四种场景。

## 2. 依赖关系

### 标准库依赖
- `encoding/json` - JSON 类型定义
- `testing` - Go 测试框架

### 第三方依赖
- `github.com/stretchr/testify/require` - 测试断言库

### 项目内部依赖
- `common` 包 - 被测试的 `JsonRawMessageToString` 函数

## 3. 类型定义

### 测试用例结构体
```go
type struct {
    name string           // 测试用例名称
    data json.RawMessage  // 输入数据
    want string           // 期望输出
}
```

## 4. 函数详解

### TestJsonRawMessageToString(t *testing.T)
```go
func TestJsonRawMessageToString(t *testing.T)
```
测试 `JsonRawMessageToString` 函数的各种输入场景。

**测试用例：**

| 用例名称 | 输入数据 | 期望输出 | 说明 |
|---------|---------|---------|------|
| object | `{"city":"Paris","days":0,"strict":false}` | `{"city":"Paris","days":0,"strict":false}` | 普通 JSON 对象 |
| string | `"{\"city\":\"Paris\",\"days\":0,\"strict\":false}"` | `{"city":"Paris","days":0,"strict":false}` | JSON 字符串解码 |
| null | `null` | `""` | null 值处理 |
| empty | `nil` | `""` | 空值处理 |

## 5. 关键逻辑分析

### 测试设计模式
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        require.Equal(t, tt.want, JsonRawMessageToString(tt.data))
    })
}
```
- 使用表驱动测试（Table-Driven Tests）模式
- 每个测试用例独立运行，通过 `t.Run` 实现子测试
- 使用 `require.Equal` 进行断言，失败时立即终止

### 测试覆盖范围
- **正常路径**：有效 JSON 对象和字符串
- **边界情况**：null 值和 nil 输入
- **类型转换**：JSON 字符串类型的自动解码

## 6. 关联文件

- `common/json.go` - 被测试的源文件，包含 `JsonRawMessageToString` 函数
- `common/json_bench_test.go` - 如果存在，可能包含性能测试
