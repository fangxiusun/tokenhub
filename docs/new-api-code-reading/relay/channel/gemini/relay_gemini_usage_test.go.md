# relay_gemini_usage_test.go 代码阅读文档

## 1. 全局总结
relay_gemini_usage_test.go 测试 Gemini 适配器的使用量提取功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `relay/channel/gemini/relay_gemini.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestExtractUsage` - 测试使用量提取
- `TestCalculateTokenCount` - 测试 Token 计数

## 4. 函数详解
### 4.1 TestExtractUsage
- **职责**: 测试从 Gemini 响应中提取使用量
- **验证点**: promptTokenCount, candidatesTokenCount, totalTokenCount

## 5. 关键逻辑分析
- **使用量格式**: 测试 Gemini 特定的使用量格式

## 6. 关联文件
- `relay/channel/gemini/relay_gemini.go` - 被测试的文件
