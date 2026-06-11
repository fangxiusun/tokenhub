# message_delta_usage_patch_test.go 代码阅读文档

## 1. 全局总结
message_delta_usage_patch_test.go 测试 Claude 消息增量使用量补丁功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `relay/channel/claude/` - Claude 适配器包

## 3. 类型定义
### 3.1 测试函数
- `TestMessageDeltaUsagePatch` - 测试增量使用量补丁

## 4. 函数详解
### 4.1 TestMessageDeltaUsagePatch
- **职责**: 测试 Claude 流式响应中的使用量提取
- **验证点**: input_tokens, output_tokens 计算

## 5. 关键逻辑分析
- **流式处理**: 测试 SSE 流中的使用量累积

## 6. 关联文件
- `relay/channel/claude/relay_claude.go` - Claude 适配器
