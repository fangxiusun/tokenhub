# compile.go 代码阅读文档

## 1. 全局总结
该文件实现了计费表达式的编译和缓存机制。它使用 expr-lang/expr 库将表达式字符串编译为可执行程序，并维护一个基于 SHA-256 哈希的缓存，以避免重复编译。同时提供版本解析、变量提取和缓存管理功能。

## 2. 依赖关系
- **fmt**: 格式化错误信息。
- **math**: 数学函数（Max、Min 等），用于编译环境。
- **strings**: 字符串处理，如前缀检查。
- **sync**: 互斥锁，保护缓存并发访问。
- **github.com/expr-lang/expr**: 表达式编译和执行库。
- **github.com/expr-lang/expr/ast**: 抽象语法树，用于变量提取。
- **github.com/expr-lang/expr/vm**: 虚拟机程序类型。

## 3. 类型定义
### cachedEntry
结构体，表示缓存条目：
- **prog**: 编译后的虚拟机程序指针。
- **usedVars**: 表达式中使用的变量集合。
- **version**: 表达式版本。

## 4. 函数详解
### ParseExprVersion
解析表达式字符串的版本标签和主体。格式为 "v1:tier(...)"，返回版本号和主体。无前缀时返回默认版本 1。

### getCompileEnv
根据版本号返回编译环境原型。当前仅支持 v1，返回 compileEnvPrototypeV1。

### CompileFromCache
编译表达式字符串，使用缓存。计算表达式的 SHA-256 哈希作为缓存键，返回编译后的程序。

### CompileFromCacheByHash
类似 CompileFromCache，但接受预计算的哈希值，用于调用方已持有 BillingSnapshot.ExprHash 的场景。

### compileFromCacheByHash
内部函数，实现编译和缓存逻辑。先检查缓存，未命中则编译表达式，提取使用的变量，并存储到缓存（缓存满时清空重建）。

### ExprVersion
返回缓存中表达式的版本。若未编译或为空，返回默认版本。

### extractUsedVars
从编译后的程序中提取使用的变量名集合，通过遍历 AST 实现。

### UsedVars
返回表达式引用的变量名集合。若未编译则先编译并缓存。

### InvalidateCache
清空编译表达式缓存，在计费规则更新时调用。

## 5. 关键逻辑分析
- **缓存策略**: 使用 SHA-256 哈希作为键，缓存编译后的程序和变量集合。缓存大小限制为 256，超出时清空重建（简单但有效）。
- **版本管理**: 支持表达式版本前缀（如 "v1:"），为未来扩展预留。
- **变量提取**: 通过 AST 遍历提取标识符节点，用于后续执行时仅传递必要变量。
- **并发安全**: 使用 RWMutex 保护缓存的并发读写。

## 6. 关联文件
- **run.go**: 调用 CompileFromCache 编译表达式，然后执行。
- **types.go**: 定义 ExprHashString 函数，用于计算表达式哈希。
- **settle.go**: 在结算时调用 CompileFromCacheByHash 编译表达式。