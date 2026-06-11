# sensitive.go 代码阅读文档

## 1. 全局摘要

该文件定义了敏感词检测 API 的响应结构 `SensitiveResponse`，用于返回检测到的敏感词列表。

## 2. 依赖

无外部包依赖。

## 3. 类型定义

### SensitiveResponse 结构体
敏感词检测响应：
- `SensitiveWords` ([]string)：敏感词列表
- `Content` (string)：原始内容

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **简单结构**：仅包含敏感词列表和原始内容。

2. **检测结果**：`SensitiveWords` 字段包含所有检测到的敏感词。

## 6. 相关文件

- `relay/sensitive/`：敏感词检测中继适配器
- `controller/sensitive.go`：敏感词检测控制器