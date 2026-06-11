# main.go 代码阅读文档

## 1. 全局总结
main.go 是应用程序的入口点，负责初始化所有组件、设置 HTTP 服务器并启动后台任务。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `constant` - 常量定义
- `i18n` - 国际化
- `logger` - 日志系统
- `middleware` - 中间件
- `model` - 数据模型
- `oauth` - OAuth 提供商
- `router` - 路由
- `service` - 服务层
- `setting` - 配置管理
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- 无（程序入口）

## 3. 类型定义
### 3.1 函数
- `main()` - 程序入口
- `InitResources()` - 初始化资源
- `setupHttpServer()` - 设置 HTTP 服务器
- `startBackgroundTasks()` - 启动后台任务

## 4. 函数详解
### 4.1 main
- **职责**: 程序入口点
- **逻辑流程**:
  1. 调用 InitResources() 初始化
  2. 调用 setupHttpServer() 设置服务器
  3. 调用 startBackgroundTasks() 启动后台任务
  4. 启动 HTTP 服务

### 4.2 InitResources
- **职责**: 初始化所有子系统
- **逻辑流程**:
  1. 初始化日志系统
  2. 初始化数据库连接
  3. 初始化 Redis
  4. 初始化 i18n
  5. 初始化 OAuth 提供商
  6. 初始化配置

### 4.3 setupHttpServer
- **职责**: 配置 Gin HTTP 服务器
- **逻辑流程**:
  1. 创建 Gin 引擎
  2. 应用中间件链
  3. 设置路由
  4. 嵌入前端资源

### 4.4 startBackgroundTasks
- **职责**: 启动维护任务
- **包含任务**:
  - 订阅过期检查
  - 渠道余额更新
  - 渠道模型同步
  - 日志清理
  - 性能指标收集

## 5. 关键逻辑分析
- **初始化顺序**: 日志 → 数据库 → Redis → i18n → OAuth → 配置
- **优雅降级**: 非关键服务（i18n、OAuth）失败不阻止启动
- **嵌入式前端**: 通过 Go embed 指令打包前端构建产物

## 6. 关联文件
- `router/main.go` - 路由初始化
- `model/main.go` - 数据库初始化
- `common/init.go` - 公共工具初始化
