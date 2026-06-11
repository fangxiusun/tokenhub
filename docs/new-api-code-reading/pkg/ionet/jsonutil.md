# jsonutil.go 代码阅读文档

## 1. 全局总结
该文件提供 JSON 解析工具函数，用于处理 API 响应中不规范的时间戳格式。支持灵活的时间戳解析，将缺少时区信息的时间字符串转换为 RFC3339Nano 格式。

## 2. 依赖关系
- **encoding/json**: JSON 序列化和反序列化。
- **strings**: 字符串处理。
- **time**: 时间解析和格式化。
- **github.com/samber/lo**: 集合操作（MapValues、Map）。

## 3. 类型定义
无类型定义。

## 4. 函数详解
### decodeWithFlexibleTimes
反序列化 JSON 数据，容忍缺少时区信息的时间戳字符串，将其标准化为 RFC3339Nano。

### decodeData
泛型函数，解包 {"data": T} 格式的响应，提取 data 字段。

### decodeDataWithFlexibleTimes
结合 decodeData 和 decodeWithFlexibleTimes，解包 data 字段并处理灵活时间戳。

### normalizeTimeValues
递归遍历 JSON 值，标准化所有时间字符串。

### normalizeTimeString
标准化单个时间字符串：尝试解析为 RFC3339Nano、RFC3339，然后尝试多种无时区格式，转换为 UTC RFC3339Nano。

## 5. 关键逻辑分析
- **时间标准化**: 处理多种时间格式，确保时区信息一致。
- **递归处理**: 递归遍历 JSON 对象和数组，处理嵌套的时间戳。
- **容错性**: 无法解析的时间字符串保持原样。

## 6. 关联文件
- **container.go、deployment.go、hardware.go**: 使用这些函数解析 API 响应。