# midjourney.go 代码阅读文档

## 1. 全局概述

本文件定义了 Midjourney 相关的常量，包括错误码、操作类型以及模型到操作的映射表。Midjourney 是一个 AI 图像生成服务，本文件为其代理集成提供标准化的常量定义。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 错误码常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `MjErrorUnknown` | `5` | 未知错误 |
| `MjRequestError` | `4` | 请求错误 |

### 操作类型常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `MjActionImagine` | `"IMAGINE"` | 文生图 |
| `MjActionDescribe` | `"DESCRIBE"` | 图片描述 |
| `MjActionBlend` | `"BLEND"` | 图片混合 |
| `MjActionUpscale` | `"UPSCALE"` | 放大图片 |
| `MjActionVariation` | `"VARIATION"` | 变体生成 |
| `MjActionReRoll` | `"REROLL"` | 重新生成 |
| `MjActionInPaint` | `"INPAINT"` | 局部重绘 |
| `MjActionModal` | `"MODAL"` | 弹窗操作 |
| `MjActionZoom` | `"ZOOM"` | 缩放 |
| `MjActionCustomZoom` | `"CUSTOM_ZOOM"` | 自定义缩放 |
| `MjActionShorten` | `"SHORTEN"` | 缩短提示词 |
| `MjActionHighVariation` | `"HIGH_VARIATION"` | 高变化度变体 |
| `MjActionLowVariation` | `"LOW_VARIATION"` | 低变化度变体 |
| `MjActionPan` | `"PAN"` | 平移 |
| `MjActionSwapFace` | `"SWAP_FACE"` | 换脸 |
| `MjActionUpload` | `"UPLOAD"` | 上传 |
| `MjActionVideo` | `"VIDEO"` | 视频 |
| `MjActionEdits` | `"EDITS"` | 编辑 |

### MidjourneyModel2Action 映射

`map[string]string` 类型，将模型名称映射到对应的操作类型：

| 模型名 | 操作 |
|--------|------|
| `mj_imagine` | IMAGINE |
| `mj_describe` | DESCRIBE |
| `mj_blend` | BLEND |
| `mj_upscale` | UPSCALE |
| `mj_variation` | VARIATION |
| `mj_reroll` | REROLL |
| `mj_modal` | MODAL |
| `mj_inpaint` | INPAINT |
| `mj_zoom` | ZOOM |
| `mj_custom_zoom` | CUSTOM_ZOOM |
| `mj_shorten` | SHORTEN |
| `mj_high_variation` | HIGH_VARIATION |
| `mj_low_variation` | LOW_VARIATION |
| `mj_pan` | PAN |
| `swap_face` | SWAP_FACE |
| `mj_upload` | UPLOAD |
| `mj_video` | VIDEO |
| `mj_edits` | EDITS |

## 6. 相关文件

- `constant/task.go` — Suno 等任务平台的常量定义
- `relay/channel/midjourney/` — Midjourney 代理适配器
- `controller/midjourney.go` — Midjourney 控制器
