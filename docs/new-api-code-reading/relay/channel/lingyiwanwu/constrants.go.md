# constrants.go 代码阅读文档

## 1. 全局总结
本文件定义了零一万物（Lingyiwanwu / Yi）渠道的常量，包括渠道名称和支持的模型列表。文件名拼写为 `constrants.go`（应为 `constants.go`），可能是历史拼写错误。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` 变量
```go
var ModelList = []string{
    "yi-large", "yi-medium", "yi-vision", "yi-medium-200k",
    "yi-spark", "yi-large-rag", "yi-large-turbo",
    "yi-large-preview", "yi-large-rag-preview",
}
```
零一万物渠道支持的模型列表：
- `yi-large` — 大规模模型
- `yi-medium` — 中等规模模型
- `yi-vision` — 视觉多模态模型
- `yi-medium-200k` — 支持 200K 上下文的中等模型
- `yi-spark` — Spark 模型
- `yi-large-rag` — 大规模 RAG 模型
- `yi-large-turbo` — 大规模 Turbo 模型
- `yi-large-preview` — 大规模预览版模型
- `yi-large-rag-preview` — 大规模 RAG 预览版模型

### `ChannelName` 变量
```go
var ChannelName = "lingyiwanwu"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析

### 零一万物（零一智能）
零一万物是李开复创办的 AI 公司，其 Yi 系列模型提供多种规模和能力。本文件注册了其主要模型，覆盖通用对话、长上下文、视觉、RAG 等场景。

### 文件名拼写
文件名 `constrants.go` 应为 `constants.go`。这个拼写错误可能是历史遗留，但由于 Go 包内所有文件共享命名空间，不影响功能。

## 6. 关联文件
- 该目录下无其他 Go 文件，零一万物渠道可能使用通用的 OpenAI 兼容处理逻辑，仅需注册模型列表。
