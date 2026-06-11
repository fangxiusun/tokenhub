# relay_claude_test.go 代码阅读文档

## 1. 全局总结
relay_claude_test.go 测试 Claude 适配器的请求转换和响应处理功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `relay/channel/claude/relay_claude.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestConvertOpenAIRequest` - 测试 OpenAI 格式到 Claude 格式转换
- `TestConvertClaudeRequest` - 测试 Claude 原生格式处理

## 4. 函数详解
### 4.1 TestConvertOpenAIRequest
- **职责**: 测试请求格式转换
- **验证点**: 消息格式、工具调用、系统提示

## 5. 关键逻辑分析
- **格式兼容**: 测试 OpenAI 到 Claude 的格式转换

## 6. 关联文件
- `relay/channel/claude/relay_claude.go` - 被测试的文件
