# waffo_pancake_test.go 代码阅读文档

## 1. 全局总结
waffo_pancake_test.go 测试 Waffo-Pancake 支付集成功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `service/waffo_pancake.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestCreateCheckoutSession` - 测试创建结账会话
- `TestVerifyWebhook` - 测试 Webhook 验证

## 4. 函数详解
### 4.1 TestCreateCheckoutSession
- **职责**: 测试 Waffo-Pancake 结账会话创建
- **验证点**: 会话 ID、支付链接

### 4.2 TestVerifyWebhook
- **职责**: 测试 Webhook 签名验证
- **验证点**: 签名有效性、防篡改

## 5. 关键逻辑分析
- **支付安全**: 测试 Webhook 签名验证

## 6. 关联文件
- `service/waffo_pancake.go` - 被测试的文件
