# billing_expr_request_test.go 代码阅读文档

## 1. 全局总结

本文件是 `billing_expr_request.go` 的单元测试，测试了计费表达式请求输入的解析和构建功能。

## 2. 依赖关系

- `gin`: 测试上下文
- `dto`: 请求 DTO
- `relay/common`: RelayInfo

## 3. 测试用例详解

### `TestResolveIncomingBillingExprRequestInput`
- 从请求上下文解析计费表达式输入
- 验证 Body 和 Headers 正确提取

### `TestBuildBillingExprRequestInputFromRequest`
- 从请求对象构建计费表达式输入
- 验证 Headers 和 Body 序列化正确

## 4. 关键逻辑分析

1. **Body 提取**: 从 BodyStorage 中读取请求体
2. **Headers 合并**: 将请求头传递到计费输入中

## 5. 关联文件

- `relay/helper/billing_expr_request.go`: 被测试的源文件
