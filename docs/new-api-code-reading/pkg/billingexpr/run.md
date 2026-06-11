# run.go 代码阅读文档

## 1. 全局总结
该文件实现了计费表达式的执行逻辑。它提供 RunExpr、RunExprWithRequest 等函数，编译表达式（使用缓存）并执行，返回计算结果和追踪信息。执行环境暴露了 token 参数、时间函数、请求探测函数等，支持复杂的计费逻辑。

## 2. 依赖关系
- **fmt**: 格式化错误信息。
- **math**: 数学函数（Max、Min 等）。
- **strings**: 字符串处理。
- **time**: 时间函数，用于时区转换。
- **github.com/expr-lang/expr**: 表达式执行库。
- **github.com/expr-lang/expr/vm**: 虚拟机程序类型。
- **github.com/tidwall/gjson**: JSON 解析，用于从请求体中提取参数。

## 3. 类型定义
无新的类型定义，但使用了 billingexpr 包中的 TokenParams、RequestInput、TraceResult 等类型。

## 4. 函数详解
### RunExpr
编译并执行表达式字符串，返回计算结果（float64）、追踪结果和错误。内部调用 RunExprWithRequest。

### RunExprWithRequest
类似 RunExpr，但接受 RequestInput 参数，支持从请求头和请求体中探测信息。

### RunExprByHash
类似 RunExpr，但接受预计算的哈希值，避免重复计算 SHA-256。

### RunExprByHashWithRequest
类似 RunExprWithRequest，但接受预计算的哈希值。

### runProgram
内部函数，执行编译后的程序。构建执行环境，包括 token 参数、层级追踪函数、请求探测函数、时间函数等。执行表达式并返回结果。

### timeInZone
将时区字符串转换为 time.Time。无效时区或空字符串默认使用 UTC。

### normalizeHeaders
标准化请求头：键转小写，值去除空格，过滤空键值对。

## 5. 关键逻辑分析
- **执行环境**: 环境变量包括 p、c、len、cr、cc、cc1h、img、img_o、ai、ao 等 token 参数，以及 tier、header、param、has、hour、minute、weekday、month、day 等函数。
- **层级追踪**: tier(name, value) 函数在执行时记录匹配的层级和成本，用于后续结算。
- **请求探测**: param(path) 使用 gjson 从请求体提取参数，header(key) 从请求头提取值，has(source, substr) 检查字符串包含关系。
- **时间函数**: 支持时区转换，用于基于时间的计费规则（如夜间折扣）。
- **错误处理**: 编译错误和执行错误均返回错误信息。

## 6. 关联文件
- **compile.go**: 提供 CompileFromCache 和 CompileFromCacheByHash 函数，用于编译表达式。
- **types.go**: 定义 TokenParams、RequestInput、TraceResult 等类型。
- **settle.go**: 调用 RunExprByHashWithRequest 执行表达式进行结算。
- **billingexpr_test.go**: 测试各种执行场景。