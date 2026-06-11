# tiered_settle_test.go 代码阅读文档

## 1. 全局总结
tiered_settle_test.go 测试分层计费结算功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `service/tiered_settle.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestTieredSettlement` - 测试分层结算
- `TestExpressionBilling` - 测试表达式计费

## 4. 函数详解
### 4.1 TestTieredSettlement
- **职责**: 测试基于使用量的分层计费
- **验证点**: 阶梯价格、阈值判断

### 4.2 TestExpressionBilling
- **职责**: 测试计费表达式解析和执行
- **验证点**: 变量替换、函数调用

## 5. 关键逻辑分析
- **表达式引擎**: 测试计费表达式的灵活性

## 6. 关联文件
- `service/tiered_settle.go` - 被测试的文件
- `pkg/billingexpr/` - 计费表达式引擎
