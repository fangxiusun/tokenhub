# relay_aws_test.go 代码阅读文档

## 1. 全局总结
relay_aws_test.go 是 AWS Bedrock 适配器的测试文件，测试 AWS 签名、请求构建和响应解析等功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架
- `relay/channel/aws` - 被测试的包

### 2.2 被引用的文件
- `relay/channel/aws/relay_aws.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestAWSSigning` - 测试 AWS 签名
- `TestBuildRequestURL` - 测试请求 URL 构建
- `TestParseResponse` - 测试响应解析

## 4. 函数详解
### 4.1 TestAWSSigning
- **职责**: 测试 AWS SigV4 签名生成
- **验证点**: 签名头格式、时间戳、日期

### 4.2 TestBuildRequestURL
- **职责**: 测试 AWS Bedrock URL 构建
- **验证点**: 区域、服务、端点

## 5. 关键逻辑分析
- **签名验证**: 确保 AWS 请求签名正确
- **区域支持**: 测试不同 AWS 区域

## 6. 关联文件
- `relay/channel/aws/relay_aws.go` - 被测试的文件
