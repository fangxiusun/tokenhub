# ratio_sync.go 代码阅读文档

## 1. 全局摘要

该文件定义了比率同步功能的数据结构，包括上游信息 `UpstreamDTO`、同步请求 `UpstreamRequest`、测试结果 `TestResult`、差异项 `DifferenceItem`，以及可同步渠道 `SyncableChannel`。

## 2. 依赖

无外部包依赖。

## 3. 类型定义

### UpstreamDTO 结构体
上游信息数据传输对象：
- `ID` (int)：上游 ID
- `Name` (string)：上游名称（必填）
- `BaseURL` (string)：基础 URL（必填）
- `Endpoint` (string)：端点

### UpstreamRequest 结构体
同步请求结构：
- `ChannelIDs` ([]int64)：渠道 ID 数组
- `Upstreams` ([]UpstreamDTO)：上游数组
- `Timeout` (int)：超时时间

### TestResult 结构体
测试结果结构：
- `Name` (string)：测试名称
- `Status` (string)：测试状态
- `Error` (string)：错误信息（可选）

### DifferenceItem 结构体
差异项结构：
- `Current` (interface{})：本地当前值
- `Upstreams` (map[string]interface{})：各渠道的上游值
- `Confidence` (map[string]bool)：置信度映射

### SyncableChannel 结构体
可同步渠道结构：
- `ID` (int)：渠道 ID
- `Name` (string)：渠道名称
- `BaseURL` (string)：基础 URL
- `Status` (int)：状态
- `Type` (int)：类型

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **必填字段**：`UpstreamDTO` 使用 `binding:"required"` 标签确保必填字段。

2. **多上游支持**：`UpstreamRequest` 支持同时指定多个上游进行同步。

3. **差异比较**：`DifferenceItem` 支持本地值与多渠道上游值的比较。

4. **置信度管理**：`Confidence` 字段记录每个上游值的置信度。

## 6. 相关文件

- `controller/ratio_sync.go`：比率同步控制器
- `service/ratio_sync.go`：比率同步服务
- `model/channel.go`：渠道数据模型