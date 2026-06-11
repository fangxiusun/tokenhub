# container.go 代码阅读文档

## 1. 全局总结
该文件实现了容器相关 API 操作，包括列出容器、获取容器详情、获取容器日志、流式日志、重启/停止容器、执行命令等。支持日志分页和流式处理。

## 2. 依赖关系
- **encoding/json**: JSON 解析。
- **fmt**: 格式化错误信息。
- **strings**: 字符串处理。
- **time**: 时间处理。
- **github.com/samber/lo**: 集合操作（FilterMap）。

## 3. 类型定义
无新的类型定义，但使用了 types.go 中的 Container、ContainerList、ContainerLogs、LogEntry 等类型。

## 4. 函数详解
### ListContainers
列出指定部署的所有容器。

### GetContainerDetails
获取指定容器的详细信息。

### GetContainerJobs
获取指定容器的任务列表。

### buildLogEndpoint
构建日志请求端点，支持日志级别、流、限制、游标、跟随等选项。

### GetContainerLogs
获取容器日志并标准化为 LogEntry 列表。

### GetContainerLogsRaw
获取原始文本日志。

### StreamContainerLogs
流式获取容器日志，使用回调函数处理每个日志条目。支持分页和轮询。

### RestartContainer
重启指定容器。

### StopContainer
停止指定容器。

### ExecuteInContainer
在容器中执行命令，返回输出。

## 5. 关键逻辑分析
- **日志处理**: 支持原始日志和结构化日志，标准化换行符，过滤空行。
- **流式日志**: 使用轮询机制模拟流式日志，支持游标分页。
- **命令执行**: 发送命令数组，返回输出字符串。

## 6. 关联文件
- **client.go**: 使用 makeRequest 方法执行 API 调用。
- **jsonutil.go**: 使用 decodeWithFlexibleTimes 和 decodeDataWithFlexibleTimes 解析响应。
- **types.go**: 定义容器相关类型。