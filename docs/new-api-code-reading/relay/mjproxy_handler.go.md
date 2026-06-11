# mjproxy_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 Midjourney 代理的完整功能，包括任务提交、状态查询、图像代理、回调通知、换脸等。是 Midjourney Proxy 协议的服务端实现，处理与 Midjourney Discord bot 的交互。

## 2. 依赖关系

- `model`: Midjourney 任务模型、渠道模型
- `service`: HTTP 请求、计费、Midjourney 错误处理
- `relay/common`: RelayInfo
- `relay/helper`: 价格计算
- `setting`: Midjourney 设置

## 3. 类型定义

### `taskChangeParams`
```go
type taskChangeParams struct {
    ID     string
    Action string
    Index  int
}
```

## 4. 函数详解

### `RelayMidjourneyImage(c *gin.Context)`
- 代理 Midjourney 图像，支持 SSRF 防护和代理配置

### `RelayMidjourneyNotify(c *gin.Context) *dto.MidjourneyResponse`
- 处理 Midjourney webhook 回调，更新任务状态

### `coverMidjourneyTaskDto(c, originTask) dto.MidjourneyDto`
- 将数据库模型转换为 DTO，处理图像 URL 代理

### `RelaySwapFace(c, info) *dto.MidjourneyResponse`
- 换脸功能，需要 source 和 target 的 base64 图像

### `RelayMidjourneyTaskImageSeed(c *gin.Context) *dto.MidjourneyResponse`
- 获取任务图像的 seed 值

### `RelayMidjourneyTask(c, relayMode) *dto.MidjourneyResponse`
- 查询任务状态，支持单个和批量查询

### `RelayMidjourneySubmit(c, relayInfo) *dto.MidjourneyResponse`
- 提交 Midjourney 任务（绘画/描述/编辑/混合/上传等）
- 支持多种操作类型和 Plus 版本的 action/modal/shorten

### `getMjRequestPath(path) string`
- 从请求路径中提取 MJ API 路径

## 5. 关键逻辑分析

1. **任务状态码处理**:
   - 1: 提交成功
   - 21: 任务已存在（处理中或有结果）
   - 22: 排队中
   - 23: 队列已满
   - 24: 敏感词
   - 3: 无可用实例（自动禁用渠道）

2. **图像代理**: 通过 `/mj/image/{id}` 代理图像，支持 SSRF 防护

3. **渠道锁定**: 放大/变换等操作必须使用原始任务的渠道

4. **计费逻辑**: 使用固定价格模型，成功后才扣费

## 6. 关联文件

- `model/midjourney.go`: Midjourney 任务模型
- `dto/midjourney.go`: Midjourney DTO
- `service/midjourney.go`: Midjourney 服务
