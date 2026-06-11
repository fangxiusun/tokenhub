# outbound_body.go 代码阅读文档

## 1. 全局总结

本文件提供了出站请求体的封装函数 `NewOutboundJSONBody`，将已序列化的 JSON 数据包装为 `BodyStorage`，支持磁盘缓存以减少大 payload（如 base64 图像）的堆内存占用。

## 2. 依赖关系

- `common`: BodyStorage、ReaderOnly

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `NewOutboundJSONBody(data []byte) (body io.Reader, size int64, closer io.Closer, err error)`
- **功能**: 将 JSON 数据包装为出站请求体
- **特性**:
  - 当启用磁盘缓存且 payload 超过阈值时，写入临时文件
  - 内存模式下复用相同的底层数组
  - 返回 `ReaderOnly` 包装防止 HTTP transport 提前关闭
  - 返回 size 用于设置 `http.Request.ContentLength`

## 5. 关键逻辑分析

1. **内存优化**: 大 base64 payload 可以写入磁盘，减少堆内存压力
2. **GC 友好**: 写入磁盘后原始 `[]byte` 可以被 GC 回收
3. **必须关闭**: 调用方必须在请求完成后调用 `closer.Close()` 释放资源

## 6. 关联文件

- `common/body_storage.go`: BodyStorage 实现
