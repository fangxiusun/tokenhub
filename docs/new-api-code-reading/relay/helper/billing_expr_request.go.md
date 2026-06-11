# billing_expr_request.go 代码阅读文档

## 1. 全局总结

本文件提供了分层计费表达式系统所需的请求输入解析功能，用于从请求上下文和请求体中提取计费表达式所需的变量值。

## 2. 依赖关系

- `common`: JSON 序列化、BodyStorage
- `dto`: Request 接口
- `pkg/billingexpr`: RequestInput 类型
- `relay/common`: RelayInfo

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ResolveIncomingBillingExprRequestInput(c, info) (billingexpr.RequestInput, error)`
- 从请求上下文解析计费表达式输入
- 优先使用已缓存的 BillingRequestInput
- 否则从请求体和 headers 中提取

### `BuildBillingExprRequestInputFromRequest(request, headers) (billingexpr.RequestInput, error)`
- 从请求对象构建计费表达式输入
- 将请求序列化为 JSON 作为 body

### `readIncomingBillingExprBody(c) ([]byte, error)`
- 读取请求体，仅处理 JSON content type

## 5. 关键逻辑分析

1. **缓存优先**: 如果 info.BillingRequestInput 已存在，直接使用并合并 headers
2. **深拷贝**: 使用 append([]byte(nil), ...) 进行 body 深拷贝
3. **Content-Type 检查**: 仅处理 application/json 类型的请求体

## 6. 关联文件

- `pkg/billingexpr/expr.md`: 分层计费表达式设计文档
- `relay/helper/price.go`: modelPriceHelperTiered 使用此函数
