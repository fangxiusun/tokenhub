# task.go 代码阅读文档

## 1. 全局概述

本文件定义了异步任务平台（Task Platform）相关的常量，包括平台类型、操作类型以及 Suno 模型到操作的映射表。支持的任务平台包括 Suno（音乐生成）和 Midjourney（图像生成）。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### TaskPlatform 类型

```go
type TaskPlatform string
```

字符串类型的任务平台标识符。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 任务平台常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `TaskPlatformSuno` | `"suno"` | Suno 音乐生成平台 |
| `TaskPlatformMidjourney` | `"mj"` | Midjourney 图像生成平台 |

### Suno 操作类型

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `SunoActionMusic` | `"MUSIC"` | 音乐生成 |
| `SunoActionLyrics` | `"LYRICS"` | 歌词生成 |

### 通用任务操作类型

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `TaskActionGenerate` | `"generate"` | 通用生成 |
| `TaskActionTextGenerate` | `"textGenerate"` | 文本生成 |
| `TaskActionFirstTailGenerate` | `"firstTailGenerate"` | 首尾生成 |
| `TaskActionReferenceGenerate` | `"referenceGenerate"` | 参考生成 |
| `TaskActionRemix` | `"remixGenerate"` | 混音生成 |

### SunoModel2Action 映射

```go
var SunoModel2Action = map[string]string{
    "suno_music":  SunoActionMusic,
    "suno_lyrics": SunoActionLyrics,
}
```

将 Suno 模型名称映射到对应的操作类型。

## 6. 相关文件

- `constant/midjourney.go` — Midjourney 相关常量
- `relay/channel/suno/` — Suno 适配器
- `relay/channel/midjourney/` — Midjourney 适配器
- `controller/task.go` — 任务控制器
- `model/task.go` — 任务数据模型
