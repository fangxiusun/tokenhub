# constants.go 代码阅读文档

## 1. 全局总结

该文件定义了 MiniMax 海螺视频生成频道的所有常量，包括频道名称、模型列表、API 端点、HTTP 状态码、任务状态、分辨率选项以及默认配置。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### 字符串常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `ChannelName` | `"hailuo-video"` | 频道标识 |

### 模型列表

| 模型名 | 说明 |
|--------|------|
| `MiniMax-Hailuo-2.3` | 海螺 2.3 版本 |
| `MiniMax-Hailuo-2.3-Fast` | 海螺 2.3 快速版 |
| `MiniMax-Hailuo-02` | 海螺 02 版本 |
| `T2V-01-Director` | 文生视频 01 导演版 |
| `T2V-01` | 文生视频 01 |
| `I2V-01-Director` | 图生视频 01 导演版 |
| `I2V-01-live` | 图生视频 01 直播版 |
| `I2V-01` | 图生视频 01 |
| `S2V-01` | 主体生成视频 01 |

### API 端点

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `TextToVideoEndpoint` | `/v1/video_generation` | 视频生成端点 |
| `QueryTaskEndpoint` | `/v1/query/video_generation` | 任务查询端点 |

### HTTP 状态码

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `StatusSuccess` | 0 | 成功 |
| `StatusRateLimit` | 1002 | 速率限制 |
| `StatusAuthFailed` | 1004 | 认证失败 |
| `StatusNoBalance` | 1008 | 余额不足 |
| `StatusSensitive` | 1026 | 敏感内容 |
| `StatusParamError` | 2013 | 参数错误 |
| `StatusInvalidKey` | 2049 | 无效密钥 |

### 任务状态

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `TaskStatusPreparing` | `"Preparing"` | 准备中 |
| `TaskStatusQueueing` | `"Queueing"` | 排队中 |
| `TaskStatusProcessing` | `"Processing"` | 处理中 |
| `TaskStatusSuccess` | `"Success"` | 成功 |
| `TaskStatusFailed` | `"Fail"` | 失败 |

### 分辨率常量

| 常量名 | 值 |
|--------|-----|
| `Resolution512P` | `"512P"` |
| `Resolution720P` | `"720P"` |
| `Resolution768P` | `"768P"` |
| `Resolution1080P` | `"1080P"` |

### 默认配置

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `DefaultDuration` | 6 | 默认视频时长（秒） |
| `DefaultResolution` | `Resolution720P` | 默认分辨率 |

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

该文件为纯常量定义。状态码设计遵循海螺 API 的自定义错误码体系（非标准 HTTP 状态码），其中 0 表示成功。任务状态使用首字母大写的字符串格式。

## 6. 关联文件

- `relay/channel/task/hailuo/adaptor.go` — 引用所有常量
- `relay/channel/task/hailuo/models.go` — 使用分辨率常量定义模型配置
