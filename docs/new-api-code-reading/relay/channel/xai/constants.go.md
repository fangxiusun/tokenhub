# constants.go 代码阅读文档

## 1. 全局总结
xAI 渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    // 语言模型
    "grok-4-1-fast-reasoning",
    "grok-4-1-fast-non-reasoning",
    "grok-code-fast-1",
    "grok-4-fast-reasoning",
    "grok-4-fast-non-reasoning",
    "grok-4-0709",
    "grok-3-mini",
    "grok-3",
    "grok-2-vision-1212",
    // 搜索变体
    "grok-4-1-fast-reasoning-search",
    "grok-4-1-fast-non-reasoning-search",
    "grok-4-fast-reasoning-search",
    "grok-4-fast-non-reasoning-search",
    "grok-4-0709-search",
    "grok-3-mini-search",
    "grok-3-search",
    // 推理努力变体
    "grok-3-mini-high", "grok-3-mini-low",
    // 图片生成模型
    "grok-imagine-image-pro",
    "grok-imagine-image",
    "grok-2-image-1212",
    // 视频生成模型
    "grok-imagine-video",
}
```

### ChannelName
```go
var ChannelName = "xai"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 模型列表包含推理/非推理变体、搜索变体、图片/视频生成模型
- 搜索变体通过 `-search` 后缀标识
- 推理努力变体通过 `-high`/`-low` 后缀标识

## 6. 关联文件
- `xai/adaptor.go` — 使用 ModelList 和 ChannelName
