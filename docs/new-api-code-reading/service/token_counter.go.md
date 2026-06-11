# token_counter.go 代码阅读文档

## 1. 全局总结

该文件实现 Token 计数功能，包括图片 Token 计算、请求 Token 估算、实时流 Token 计数、以及音频 Token 计算。是预扣费额度估算的核心。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | OpenAI 模型判断 |
| `constant` | Token 计数开关 |
| `dto` | 请求结构体 |
| `file_service` | 图片配置获取 |
| `relaycommon` | RelayInfo |

## 3. 函数详解

### `getImageToken(c, fileMeta, model, stream) (int, error)`
图片 Token 计算：
- **Patch-based 模型**（4.1-mini/nano, o4-mini 等）：32x32 patch，上限 1536
- **Tile-based 模型**（4o, o1, o3 等）：512px tile
- 低分辨率：固定 baseTokens
- 非流式/不统计媒体：固定 3*baseTokens

### `EstimateRequestToken(c, meta, info) (int, error)`
请求 Token 估算：
- 文本 token（使用 tokenizer 或估算）
- 工具调用 token（8 per tool）
- 消息格式 token（3 per message）
- 文件 token（图片/音频/视频/文件）

### `CountTokenRealtime(info, request, model) (int, int, error)`
实时流 Token 计数：
- 文本事件：SessionUpdate、TranscriptionDelta、FunctionCallArgumentsDelta
- 音频事件：AudioDelta、InputAudioBufferAppend
- 工具 token：ResponseDone 时计算

### `CountTextToken(text, model) int`
文本 Token 计数：
- OpenAI 模型：使用 tiktoken-go tokenizer
- 其他模型：使用估算器

## 4. 关键逻辑分析

1. **模型分类**：不同模型使用不同的图片 token 计算公式
2. **Patch vs Tile**：新模型使用 patch，旧模型使用 tile
3. **缩放算法**：图片超过上限时按比例缩放
4. **混合计数**：OpenAI 精确计数 + 其他模型估算

## 5. 关联文件

- `tokenizer.go` — Tokenizer 实现
- `token_estimator.go` — Token 估算器
- `file_service.go` — 图片配置获取
