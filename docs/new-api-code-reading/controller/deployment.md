# deployment.go 代码阅读文档

## 1. 全局总结

该文件实现了 io.net 模型部署管理的完整 API，包括部署的 CRUD 操作、硬件类型查询、位置查询、副本数查询、价格估算、集群名称验证、容器管理和日志查看。

## 2. 依赖关系

- `common` — 通用工具、Option 读取
- `pkg/ionet` — io.net API 客户端
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义（使用内联 struct 和 map）。

## 4. 函数详解

### 配置和连接
- `GetModelDeploymentSettings` — 获取部署设置状态
- `TestIoNetConnection` — 测试 io.net API 连接

### 部署管理
- `GetAllDeployments` — 获取所有部署（分页）
- `SearchDeployments` — 搜索部署
- `GetDeployment` — 获取单个部署详情
- `CreateDeployment` — 创建新部署
- `UpdateDeployment` — 更新部署
- `UpdateDeploymentName` — 更新部署名称（含可用性检查）
- `ExtendDeployment` — 延长部署时间
- `DeleteDeployment` — 终止部署

### 资源查询
- `GetHardwareTypes` — 获取可用硬件类型
- `GetLocations` — 获取可用位置
- `GetAvailableReplicas` — 获取可用副本数
- `GetPriceEstimation` — 价格估算
- `CheckClusterNameAvailability` — 集群名称可用性检查

### 容器管理
- `ListDeploymentContainers` — 列出部署中的容器
- `GetContainerDetails` — 获取容器详情
- `GetDeploymentLogs` — 获取容器日志

## 5. 关键逻辑分析

- API Key 从 `model_deployment.ionet.api_key` Option 读取
- 支持普通客户端和企业客户端两种模式
- 部署状态统计：running、completed、failed、deployment requested、termination requested、destroyed
- 日志支持分页、时间范围、级别过滤和流式跟随

## 6. 关联文件

- `pkg/ionet/` — io.net API 客户端实现
