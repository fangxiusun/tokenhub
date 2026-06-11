# error_test.go 代码阅读文档

## 1. 全局总结
error_test.go 测试服务层错误处理功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `service/error.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestRelayErrorHandler` - 测试中继错误处理器

## 4. 函数详解
### 4.1 TestRelayErrorHandler
- **职责**: 测试上游错误的解析和标准化
- **验证点**: 错误消息提取、状态码映射

## 5. 关键逻辑分析
- **错误标准化**: 测试不同提供商的错误格式统一

## 6. 关联文件
- `service/error.go` - 被测试的文件
