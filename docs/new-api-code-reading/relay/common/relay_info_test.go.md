# relay_info_test.go 代码阅读文档

## 1. 全局总结

本文件是 `relay_info.go` 的单元测试，主要测试 `GetFinalRequestRelayFormat` 方法的优先级逻辑。

## 2. 依赖关系

- `types`: RelayFormat 类型

## 3. 测试用例详解

### `TestRelayInfoGetFinalRequestRelayFormatPrefersExplicitFinal`
- 当 FinalRequestRelayFormat 被显式设置时，优先使用它

### `TestRelayInfoGetFinalRequestRelayFormatFallsBackToConversionChain`
- 未设置 FinalRequestRelayFormat 时，使用 RequestConversionChain 的最后一项

### `TestRelayInfoGetFinalRequestRelayFormatFallsBackToRelayFormat`
- 未设置 FinalRequestRelayFormat 和 RequestConversionChain 时，使用 RelayFormat

### `TestRelayInfoGetFinalRequestRelayFormatNilReceiver`
- nil 指针安全

## 4. 关键逻辑分析

1. **优先级**: FinalRequestRelayFormat > RequestConversionChain 最后一项 > RelayFormat
2. **nil 安全**: nil 指针返回空字符串

## 5. 关联文件

- `relay/common/relay_info.go`: GetFinalRequestRelayFormat 方法
