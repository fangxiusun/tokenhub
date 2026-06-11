# sensitive.go 代码阅读文档

## 1. 全局总结

该文件实现敏感词检测和过滤功能，支持消息内容检查、文本检查、以及敏感词替换。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `dto` | Message 结构体 |
| `setting` | 敏感词列表 |

## 3. 函数详解

### `CheckSensitiveMessages(messages []dto.Message) ([]string, error)`
检查消息列表中的敏感词

### `CheckSensitiveText(text string) (bool, []string)`
检查文本是否包含敏感词

### `SensitiveWordContains(text string) (bool, []string)`
使用 Aho-Corasick 自动机进行多模式匹配

### `SensitiveWordReplace(text, returnImmediately) (bool, []string, string)`
敏感词替换：
- 将敏感词替换为 `**###**`
- 返回替换后的文本

## 4. 关键逻辑分析

1. **AC 自动机**：使用 Aho-Corasick 算法进行高效的多模式匹配
2. **大小写不敏感**：统一转换为小写后匹配
3. **智能缓存**：AC 机器按字典哈希缓存，避免重复构建

## 5. 关联文件

- `str.go` — AC 自动机实现
- `setting` — 敏感词配置
