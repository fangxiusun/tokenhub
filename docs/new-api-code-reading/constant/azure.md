# azure.go 代码阅读文档

## 1. 全局概述

本文件定义了一个与 Azure 相关的时间戳常量 `AzureNoRemoveDotTime`，用于 Azure API 版本兼容性处理。

## 2. 依赖关系

- `time` — Go 标准库时间包

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### AzureNoRemoveDotTime 变量

```go
var AzureNoRemoveDotTime = time.Date(2025, time.May, 10, 0, 0, 0, 0, time.UTC).Unix()
```

该变量存储了 2025 年 5 月 10 日的 Unix 时间戳（秒级）。根据命名推断，这是一个分界点时间：
- 在该时间之前创建的 Azure 渠道，可能需要移除 API 版本号中的点号（如 `2024-02-15` → `20240215`）
- 在该时间之后创建的渠道，则不需要此处理

这是一个临时的兼容性处理方案，用于平滑过渡 Azure API 版本格式变更。

## 6. 相关文件

- `constant/channel.go` — 定义 `ChannelTypeAzure` 常量
- `relay/channel/azure/` — Azure 渠道适配器实现
