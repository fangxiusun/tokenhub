# types.go 代码阅读文档

## 1. 全局总结
该文件定义了 IO.NET API 客户端的所有数据类型，包括客户端、HTTP 请求/响应、部署、容器、硬件、位置、价格估算、日志等结构体。提供完整的 API 交互类型定义。

## 2. 依赖关系
- **time**: 时间类型。

## 3. 类型定义
### Client
结构体，API 客户端：BaseURL、APIKey、HTTPClient。

### HTTPClient
接口，定义 Do 方法。

### HTTPRequest
结构体，HTTP 请求：Method、URL、Headers、Body。

### HTTPResponse
结构体，HTTP 响应：StatusCode、Headers、Body。

### DeploymentRequest
结构体，部署请求：资源名称、时长、GPU 配置、位置、容器配置、注册表配置。

### ContainerConfig
结构体，容器配置：副本数、环境变量、入口点、流量端口等。

### RegistryConfig
结构体，注册表配置：镜像 URL、用户名、密钥。

### DeploymentResponse
结构体，部署响应：部署 ID、状态。

### DeploymentDetail
结构体，部署详情：ID、状态、时间、成本、GPU 信息、位置、容器配置等。

### Container
结构体，容器：设备 ID、容器 ID、硬件、状态、事件等。

### ContainerList
结构体，容器列表：总数、容器数组。

### Deployment
结构体，部署列表项：ID、状态、名称、硬件数量等。

### DeploymentList
结构体，部署列表：部署数组、总数、状态。

### AvailableReplica
结构体，可用副本：位置、硬件、可用数量。

### MaxGPUResponse
结构体，最大 GPU 响应：硬件数组、总数。

### PriceEstimationRequest
结构体，价格估算请求：位置、硬件、GPU 配置、时长、货币等。

### PriceEstimationResponse
结构体，价格估算响应：估算成本、货币、价格明细。

### PriceBreakdown
结构体，价格明细：计算成本、网络成本、存储成本、总成本、每小时费率。

### ContainerLogs
结构体，容器日志：容器 ID、日志条目、是否有更多、下一个游标。

### LogEntry
结构体，日志条目：时间戳、级别、消息、来源。

### UpdateDeploymentRequest
结构体，更新部署请求：环境变量、入口点、流量端口、镜像 URL 等。

### ExtendDurationRequest
结构体，扩展时长请求：时长小时数。

### UpdateDeploymentResponse
结构体，更新部署响应：状态、部署 ID。

### UpdateClusterNameRequest
结构体，更新集群名称请求：名称。

### UpdateClusterNameResponse
结构体，更新集群名称响应：状态、消息。

### APIError
结构体，API 错误：代码、消息、详情。实现 error 接口。

### ListDeploymentsOptions
结构体，列出部署选项：状态、位置、分页、排序。

### GetLogsOptions
结构体，获取日志选项：时间范围、级别、流、限制、游标、跟随。

### HardwareType
结构体，硬件类型：ID、名称、GPU 类型、内存、最大 GPU 数、时长费率等。

### Location
结构体，位置：ID、名称、ISO2、区域、国家、经纬度、可用数量。

### LocationsResponse
结构体，位置响应：位置数组、总数。

### LocationAvailability
结构体，位置可用性：位置 ID、名称、可用性、硬件可用性、更新时间。

### HardwareAvailability
结构体，硬件可用性：硬件 ID、名称、可用数量、最大 GPU 数。

## 4. 函数详解
### APIError.Error
实现 error 接口，返回错误消息（包含详情）。

## 5. 关键逻辑分析
- **类型设计**: 类型设计考虑了 API 响应格式、可选字段、派生字段。
- **JSON 标签**: 使用 json 标签定义字段名和 omitempty 选项。
- **错误处理**: APIError 实现 error 接口，便于错误处理。

## 6. 关联文件
- **client.go、container.go、deployment.go、hardware.go**: 使用这些类型进行 API 交互。