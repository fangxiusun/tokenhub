# billing.go 代码阅读文档

## 1. 全局总结

该文件包含 Gemini Veo 视频生成的计费相关工具函数，负责从元数据和标准请求参数中解析时长、分辨率，并提供尺寸到 Veo 分辨率/宽高比的转换、以及分辨率计费倍率计算。

## 2. 依赖关系

**标准库：**
- `strconv` — 字符串到整数转换
- `strings` — 字符串处理

无项目内部或第三方依赖。

## 3. 类型定义

无类型定义。

## 4. 函数详解

| 函数签名 | 说明 |
|----------|------|
| `ParseVeoDurationSeconds(metadata map[string]any) int` | 从 metadata 中提取 `durationSeconds`，默认返回 8 |
| `ParseVeoResolution(metadata map[string]any) string` | 从 metadata 中提取 `resolution`，默认返回 "720p" |
| `ResolveVeoDuration(metadata, stdDuration, stdSeconds) int` | 按优先级解析时长：metadata > stdDuration > stdSeconds > 默认 8 |
| `ResolveVeoResolution(metadata, stdSize) string` | 按优先级解析分辨率：metadata > SizeToVeoResolution(stdSize) > 默认 "720p" |
| `SizeToVeoResolution(size string) string` | 将 "WxH" 尺寸转换为 Veo 分辨率标签（720p/1080p/4k） |
| `SizeToVeoAspectRatio(size string) string` | 将 "WxH" 尺寸转换为宽高比（16:9 或 9:16） |
| `VeoResolutionRatio(modelName, resolution string) float64` | 返回分辨率计费倍率，4K 按模型返回不同倍率 |

## 5. 关键逻辑分析

### 时长解析优先级
```
metadata["durationSeconds"] > req.Duration > req.Seconds > 8（默认值）
```

### 分辨率解析优先级
```
metadata["resolution"] > SizeToVeoResolution(req.Size) > "720p"（默认值）
```

### 尺寸到分辨率转换
- 最大维度 >= 3840 → "4k"
- 最大维度 >= 1920 → "1080p"
- 其他 → "720p"

### 宽高比判断
- 高度 > 宽度 → "9:16"（竖屏）
- 其他 → "16:9"（横屏）

### 4K 计费倍率（基于 Vertex AI 官方定价）
- veo-3.1-fast-generate: 2.333（$0.35 / $0.15）
- veo-3.1-generate / veo-3.1: 1.5（$0.60 / $0.40）
- 其他模型（不支持 4K）: 1.0

## 6. 关联文件

- `relay/channel/task/gemini/adaptor.go` — 调用 `ResolveVeoDuration`、`ResolveVeoResolution`、`VeoResolutionRatio` 进行计费估算和请求构建
