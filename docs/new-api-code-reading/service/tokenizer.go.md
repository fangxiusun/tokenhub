# tokenizer.go 代码阅读文档

## 1. 全局总结

该文件封装 tiktoken-go 库，提供 OpenAI 模型的精确 Token 计数功能。使用缓存的 codec 实例避免重复创建。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `tiktoken-go/tokenizer` | Token 编解码器 |
| `tiktoken-go/tokenizer/codec` | 默认编码器 |

## 3. 类型定义

全局变量：
- `defaultTokenEncoder` — 默认编码器（cl100k_base）
- `tokenEncoderMap` — 按模型缓存的编码器
- `tokenEncoderMutex` — 并发锁

## 4. 函数详解

### `InitTokenEncoders()`
初始化默认编码器

### `getTokenEncoder(model) tokenizer.Codec`
获取模型对应的编码器：
1. 读锁检查缓存
2. 写锁创建新编码器
3. 双重检查防止重复创建
4. 失败时缓存默认编码器

### `getTokenNum(tokenEncoder, text) int`
使用编码器计算 token 数

## 5. 关键逻辑分析

1. **懒加载**：首次使用时创建编码器
2. **双重检查锁**：防止并发重复创建
3. **失败回退**：未知模型使用默认编码器

## 6. 关联文件

- `token_counter.go` — 使用 tokenizer
