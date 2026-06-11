# midjourney.go 代码阅读文档

## 1. 全局总结

该文件实现 Midjourney API 的请求处理，包括动作解析、请求转发、以及 Plus 模式的自定义 ID 解析。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | JSON 解析 |
| `constant` | 动作常量 |
| `dto` | 请求/响应结构体 |
| `logger` | 日志 |
| `relayconstant` | Relay 模式常量 |
| `setting` | MJ 配置 |

## 3. 函数详解

### `CovertMjpActionToModelName(mjAction) string`
将 MJ 动作转换为模型名称（如 `mj_imagine`）

### `GetMjRequestModel(relayMode, midjRequest) (string, *dto.MidjourneyResponse, bool)`
根据 Relay 模式获取 MJ 请求模型：
- 支持：Imagine、Video、Edits、Describe、Blend、Shorten、Change、Modal、SwapFace、Upload
- Plus 模式：解析 customId

### `CoverPlusActionToNormalAction(midjRequest)`
解析 Plus 模式的 customId 为标准动作：
- `MJ::JOB::upsample::2::...` → UPSCALE
- `MJ::JOB::variation::...` → VARIATION
- 支持：upsample、variation、pan、reroll、Outpaint、CustomZoom、Inpaint

### `ConvertSimpleChangeParams(content)`
解析简单变更参数（如 `taskId u1`）

### `DoMidjourneyHttpRequest(c, timeout, fullRequestURL)`
执行 MJ HTTP 请求：
- 读取并处理请求体
- 可选移除 accountFilter/notifyHook
- 可选清除 fast/relax/turbo 模式
- 设置认证头

## 4. 关键逻辑分析

1. **动作映射**：relayMode → MJ action → 模型名称
2. **Plus 模式解析**：从 customId 中提取动作和索引
3. **模式清除**：可选移除 fast/relax/turbo 参数

## 5. 关联文件

- `relay/midjourney/` — MJ Relay 适配器
- `constant/midjourney.go` — MJ 动作常量
