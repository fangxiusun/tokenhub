# param_override_error.go 代码阅读文档

## 1. 全局总结

本文件提供了参数覆盖（Param Override）错误的转换工具函数，将底层的 `ParamOverrideReturnError` 转换为统一的 `NewAPIError` 格式。

## 2. 依赖关系

- `relay/common`: ParamOverrideReturnError 类型和转换函数
- `types`: NewAPIError 类型

## 3. 类型定义

本文件无自定义类型定义。

## 4. 类型定义

### `newAPIErrorFromParamOverride(err error) *types.NewAPIError`
- **功能**: 将参数覆盖错误转换为 API 错误
- **逻辑**: 尝试断言为 `ParamOverrideReturnError`，成功则调用 `NewAPIErrorFromParamOverride`，否则创建通用错误
- **用途**: 在所有 handler 中统一处理参数覆盖产生的错误

## 5. 关键逻辑分析

1. **错误包装**: 将参数覆盖的自定义错误（含 status_code、code、type）转换为标准 API 错误格式
2. **统一入口**: 所有 handler 中参数覆盖错误都通过此函数转换

## 6. 关联文件

- `relay/common/override.go`: ParamOverrideReturnError 定义
- `types/error.go`: NewAPIError 类型
