# request_common.go 代码阅读文档

## 1. 全局摘要

该文件定义了请求接口 `Request` 和基础请求结构 `BaseRequest`，为所有请求 DTO 提供统一的接口规范。所有请求结构体都应实现 `Request` 接口。

## 2. 依赖

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`TokenCountMeta` 类型
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### Request 接口
所有请求 DTO 必须实现的接口：
```go
type Request interface {
    GetTokenCountMeta() *types.TokenCountMeta
    IsStream(c *gin.Context) bool
    SetModelName(modelName string)
}
```

### BaseRequest 结构体
基础请求结构，提供接口的默认实现：
- 空结构体，作为其他请求结构体的嵌入基础

## 4. 函数详情

### GetTokenCountMeta()
```go
func (b *BaseRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取 token 计数元数据（默认实现）。

**返回**：返回空的元数据，token 类型为 `TokenTypeTokenizer`。

### IsStream()
```go
func (b *BaseRequest) IsStream(c *gin.Context) bool
```
**功能**：判断是否为流式请求（默认实现）。

**返回**：始终返回 `false`。

### SetModelName()
```go
func (b *BaseRequest) SetModelName(modelName string)
```
**功能**：设置模型名称（默认实现，空操作）。

## 5. 关键逻辑分析

1. **接口规范**：定义所有请求 DTO 必须实现的接口，确保一致性。

2. **默认实现**：`BaseRequest` 提供接口的默认实现，其他请求结构体可以嵌入并覆盖。

3. **统一接口**：通过接口实现多态，便于在中继层统一处理不同类型的请求。

## 6. 相关文件

- `dto/audio.go`：音频请求实现
- `dto/openai_request.go`：OpenAI 请求实现
- `dto/claude.go`：Claude 请求实现
- `relay/handler.go`：请求处理逻辑