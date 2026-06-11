# model_mapped.go 代码阅读文档

## 1. 全局总结

本文件实现了模型名称映射逻辑，支持渠道级别的模型重定向。当渠道配置了 model_mapping 时，自动将用户请求的模型名称映射为上游渠道支持的模型名称。

## 2. 依赖关系

- `dto`: Request 接口
- `relay/common`: RelayInfo
- `relay/constant`: RelayMode
- `ratio_setting`: CompactModelSuffix

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ModelMappedHelper(c, info, request) error`
- **功能**: 执行模型名称映射
- **逻辑**:
  1. 处理 ResponsesCompact 模式的后缀
  2. 解析渠道的 model_mapping JSON
  3. 支持链式重定向（A → B → C）
  4. 循环检测（避免无限循环）
  5. 更新 info.UpstreamModelName 和 request 的模型名

## 5. 关键逻辑分析

1. **链式重定向**: 支持多级模型映射（如 gpt-4 → gpt-4-turbo → gpt-4-turbo-2024-04-09）
2. **循环检测**: 使用 visitedModels map 检测循环，避免无限重定向
3. **自映射处理**: 如果映射结果与原始名称相同且是直接映射，标记为未映射
4. **CompactModelSuffix**: ResponsesCompact 模式下自动添加/移除后缀

## 6. 关联文件

- `relay/relay_adaptor.go`: GetAdaptor 使用映射后的模型名
- `dto/request.go`: Request 接口的 SetModelName 方法
