# channel.go (controller) 代码阅读文档

## 1. 全局总结
channel.go 是渠道管理控制器，处理渠道的 CRUD 操作、模型获取、批量操作、Ollama 管理等。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `dto` - 数据传输对象
- `i18n` - 国际化
- `model` - 数据模型
- `service` - 业务逻辑
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/api-router.go` - 管理 API 路由

## 3. 类型定义
### 3.1 函数
- `GetAllChannels` - 获取所有渠道
- `SearchChannels` - 搜索渠道
- `GetChannel` - 获取单个渠道
- `AddChannel` - 添加渠道
- `UpdateChannel` - 更新渠道
- `DeleteChannel` - 删除渠道
- `FetchModels` - 获取渠道模型
- `CopyChannel` - 复制渠道
- `OllamaPullModel` - Ollama 拉取模型
- `OllamaDeleteModel` - Ollama 删除模型

## 4. 函数详解
### 4.1 AddChannel
- **职责**: 创建新渠道
- **逻辑流程**:
  1. 解析请求体
  2. 验证渠道配置
  3. 创建渠道记录
  4. 创建渠道能力

### 4.2 FetchModels
- **职责**: 从上游获取可用模型列表
- **逻辑流程**:
  1. 根据渠道类型选择适配器
  2. 调用适配器获取模型列表
  3. 返回模型列表

## 5. 关键逻辑分析
- **批量操作**: 支持批量启用/禁用/删除渠道
- **标签管理**: 支持按标签批量操作
- **模型同步**: 自动从上游同步模型列表

## 6. 关联文件
- `model/channel.go` - 渠道模型
- `service/channel.go` - 渠道缓存
- `relay/relay_adaptor.go` - 适配器工厂
