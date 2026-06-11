# models.go 代码阅读文档

## 1. 全局总结

本文件定义了 Suno 音乐生成通道的常量配置，包括支持的模型列表和通道名称标识。Suno 是一个 AI 音乐生成服务，本文件为该通道提供基础元数据。

## 2. 依赖关系

无外部依赖，纯定义文件。

## 3. 类型定义

### 变量

| 名称 | 类型 | 值 | 说明 |
|------|------|-----|------|
| `ModelList` | `[]string` | `["suno_music", "suno_lyrics"]` | Suno 通道支持的模型列表 |
| `ChannelName` | `string` | `"suno"` | 通道名称标识 |

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

- `suno_music` — 音乐生成模型，用于根据提示词生成完整音乐作品。
- `suno_lyrics` — 歌词生成模型，用于根据提示词生成歌词文本。
- 两个模型分别对应 Suno API 的不同功能端点（`/submit/generate` 和 `/submit/lyrics`）。
- 该文件与 Sora 的 `constants.go` 模式一致，保持项目结构的统一性。

## 6. 关联文件

- `relay/channel/task/suno/adaptor.go` — Suno 适配器实现，引用 `ModelList` 和 `ChannelName`
- `constant/suno.go` — Suno 动作常量定义（`SunoActionMusic`、`SunoActionLyrics`）
