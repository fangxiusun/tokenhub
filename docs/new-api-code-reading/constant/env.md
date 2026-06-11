# env.go 代码阅读文档

## 1. 全局概述

本文件定义了系统运行时的环境变量/配置变量，这些变量通过启动参数或环境变量设置，控制系统的各种行为参数。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

本文件无自定义类型定义，所有变量为 Go 基础类型。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 环境变量分类

#### 流式传输相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `StreamingTimeout` | `int` | 流式传输超时时间 |
| `ForceStreamOption` | `bool` | 是否强制启用 StreamOptions |
| `StreamScannerMaxBufferMB` | `int` | 流式扫描器最大缓冲区大小（MB） |

#### Token 计数相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `CountToken` | `bool` | 是否启用 Token 计数 |
| `GetMediaToken` | `bool` | 是否获取媒体 Token |
| `GetMediaTokenNotStream` | `bool` | 非流式时是否获取媒体 Token |

#### 文件处理相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `MaxFileDownloadMB` | `int` | 最大文件下载大小（MB） |
| `MaxRequestBodyMB` | `int` | 最大请求体大小（MB） |
| `AnonymousRequestBodyLimitKB` | `int` | 匿名用户请求体限制（KB） |

#### 任务相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `UpdateTask` | `bool` | 是否启用任务更新 |
| `TaskQueryLimit` | `int` | 任务查询限制 |
| `TaskTimeoutMinutes` | `int` | 任务超时时间（分钟） |
| `TaskPricePatches` | `[]string` | Sora 任务价格补丁（临时方案） |

#### Azure 相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `AzureDefaultAPIVersion` | `string` | Azure 默认 API 版本 |

#### 通知相关
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `NotifyLimitCount` | `int` | 通知限制次数 |
| `NotificationLimitDurationMinute` | `int` | 通知限制时间窗口（分钟） |

#### 调试与其他
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `DifyDebug` | `bool` | Dify 调试模式 |
| `ErrorLogEnabled` | `bool` | 错误日志是否启用 |
| `GenerateDefaultToken` | `bool` | 是否自动生成默认 Token |
| `TrustedRedirectDomains` | `[]string` | 可信重定向域名列表（支持子域名匹配） |

### 设计说明

- 所有变量为包级全局变量，便于在系统各处直接访问
- 使用指针或直接值类型，便于判断是否设置了默认值
- `TaskPricePatches` 被标记为临时变量，将在未来移除
- `TrustedRedirectDomains` 支持子域名匹配（如 `example.com` 匹配 `sub.example.com`）

## 6. 相关文件

- `main.go` — 启动时解析环境变量并设置这些值
- `common/env.go` — 环境变量读取工具函数
- `setting/` — 配置管理模块
