# swag_video.go 代码阅读文档

## 1. 全局总结
swag_video.go 处理视频生成相关的 API 端点，包括视频生成请求和状态查询。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `dto` - 数据传输对象
- `model` - 数据模型
- `service` - 业务逻辑
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/video-router.go` - 视频路由

## 3. 类型定义
### 3.1 函数
- `CreateVideoGeneration` - 创建视频生成请求
- `GetVideoGeneration` - 获取视频生成状态

## 4. 函数详解
### 4.1 CreateVideoGeneration
- **职责**: 创建视频生成任务
- **逻辑流程**:
  1. 解析请求体
  2. 验证参数
  3. 创建任务
  4. 返回任务 ID

## 5. 关键逻辑分析
- **异步处理**: 视频生成是异步任务
- **状态轮询**: 客户端通过轮询获取生成状态

## 6. 关联文件
- `model/task.go` - 任务模型
- `service/task.go` - 任务服务
