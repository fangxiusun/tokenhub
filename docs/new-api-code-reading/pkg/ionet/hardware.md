# hardware.go 代码阅读文档

## 1. 全局总结
该文件实现了硬件和位置相关 API 操作，包括获取可用副本、获取最大 GPU 数、列出硬件类型、列出位置、获取硬件详情、获取位置详情、获取位置可用性等。

## 2. 依赖关系
- **encoding/json**: JSON 解析。
- **fmt**: 格式化错误信息。
- **strings**: 字符串处理。
- **github.com/samber/lo**: 集合操作（Map、SumBy）。

## 3. 类型定义
无新的类型定义，但使用了 types.go 中的 HardwareType、Location、AvailableReplica 等类型。

## 4. 函数详解
### GetAvailableReplicas
获取指定硬件在每个位置的可用副本数。

### GetMaxGPUsPerContainer
获取每种硬件类型的最大 GPU 数。

### ListHardwareTypes
列出可用硬件类型，从最大 GPU 端点获取数据并映射。

### ListLocations
列出可用部署位置。

### GetHardwareType
获取指定硬件类型的详情。

### GetLocation
获取指定位置的详情。

### GetLocationAvailability
获取指定位置的实时可用性。

## 5. 关键逻辑分析
- **数据聚合**: 从多个端点聚合数据，计算总可用数量。
- **数据标准化**: 标准化位置 ISO2 代码为大写。
- **错误处理**: 验证输入参数，处理 API 错误。

## 6. 关联文件
- **client.go**: 使用 makeRequest 方法执行 API 调用。
- **jsonutil.go**: 使用 decodeData 解析响应。
- **types.go**: 定义硬件和位置相关类型。