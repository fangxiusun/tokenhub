# regex.go 代码阅读文档

## 1. 全局总结

本文件提供一个带缓存的正则匹配工具函数 `matchAnyRegex`，用于判断给定字符串是否匹配一组正则模式中的任意一个。该函数被 `policy.go` 中的策略判断逻辑调用，支持渠道模型匹配等场景。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `regexp` | 正则表达式编译与匹配 |
| `sync` | 并发安全的缓存存储（`sync.Map`） |

## 3. 类型定义

### 全局变量

```go
var compiledRegexCache sync.Map // map[string]*regexp.Regexp
```

`sync.Map` 实现的正则编译缓存，键为正则模式字符串，值为编译后的 `*regexp.Regexp`。利用 `sync.Map` 的并发安全性，避免加锁开销。

## 4. 函数详解

### `matchAnyRegex(patterns []string, s string) bool`

**签名**: `func matchAnyRegex(patterns []string, s string) bool`

**功能**: 检查字符串 `s` 是否匹配 `patterns` 中的任一正则模式。

**参数**:
- `patterns`: 正则模式切片
- `s`: 待匹配的字符串

**逻辑**:
1. **边界检查**: 若 `patterns` 为空或 `s` 为空字符串，返回 `false`
2. **逐模式匹配**: 遍历每个模式，跳过空字符串模式
3. **缓存查找**: 通过 `compiledRegexCache.Load(pattern)` 查找已编译的正则
4. **编译与缓存**: 若缓存未命中，调用 `regexp.Compile` 编译并将结果存入缓存
5. **容错处理**: 若正则编译失败（无效模式），跳过该模式继续下一个，不会中断运行时流量
6. **匹配测试**: 使用 `MatchString` 进行匹配，任一匹配成功即返回 `true`
7. 全部不匹配则返回 `false`

## 5. 关键逻辑分析

1. **编译缓存**: 使用 `sync.Map` 缓存已编译的正则对象，避免重复编译开销。对于频繁调用的策略判断函数，这是重要的性能优化。

2. **容错设计**: 无效正则模式不会导致错误，而是静默跳过。注释明确说明了设计意图：`"Treat invalid patterns as non-matching to avoid breaking runtime traffic."`

3. **惰性编译**: 正则只在首次使用时编译，未使用的模式不会消耗资源。

## 6. 关联文件

- `new-api/service/openaicompat/policy.go` — 调用方，使用该函数进行模型名匹配
