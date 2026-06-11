# deployment.go 代码阅读文档

## 1. 全局总结
该文件实现了部署相关 API 操作，包括部署容器、列出部署、获取部署详情、更新部署、扩展部署、删除部署、价格估算、集群名称检查等。提供完整的部署生命周期管理。

## 2. 依赖关系
- **encoding/json**: JSON 解析。
- **fmt**: 格式化错误信息。
- **strings**: 字符串处理。
- **github.com/samber/lo**: 集合操作（Map、SumBy）。

## 3. 类型定义
无新的类型定义，但使用了 types.go 中的 DeploymentRequest、DeploymentResponse、DeploymentDetail、PriceEstimationRequest 等类型。

## 4. 函数详解
### DeployContainer
部署新容器，验证必需字段，发送 POST 请求。

### ListDeployments
列出部署，支持状态、位置、分页、排序等过滤选项。

### GetDeployment
获取部署详情。

### UpdateDeployment
更新部署配置（环境变量、入口点、流量端口等）。

### ExtendDeployment
扩展部署时长。

### DeleteDeployment
删除部署。

### GetPriceEstimation
计算部署的估算成本，支持不同的持续时间类型（小时、天、周、月）和货币。

### CheckClusterNameAvailability
检查集群名称是否可用。

### UpdateClusterName
更新集群名称。

## 5. 关键逻辑分析
- **价格估算**: 支持多种持续时间类型转换，计算每小时费率，解析 API 响应中的费用明细。
- **部署验证**: 验证必需字段，确保请求参数有效。
- **数据映射**: 将 API 响应映射到内部类型，添加派生字段（如 GPUCount）。

## 6. 关联文件
- **client.go**: 使用 makeRequest 方法执行 API 调用。
- **jsonutil.go**: 使用 decodeData 和 decodeDataWithFlexibleTimes 解析响应。
- **types.go**: 定义部署相关类型。