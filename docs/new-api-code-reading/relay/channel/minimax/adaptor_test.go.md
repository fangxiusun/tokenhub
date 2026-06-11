# adaptor_test.go 代码阅读文档

## 1. 全局总结
adaptor_test.go 测试 MiniMax 适配器的功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `relay/channel/minimax/adaptor.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestGetModelList` - 测试模型列表获取
- `TestGetRequestURL` - 测试请求 URL 构建

## 4. 函数详解
### 4.1 TestGetModelList
- **职责**: 测试返回的模型列表
- **验证点**: 模型名称、数量

## 5. 关键逻辑分析
- **模型支持**: 验证 MiniMax 支持的模型列表

## 6. 关联文件
- `relay/channel/minimax/adaptor.go` - 被测试的文件
