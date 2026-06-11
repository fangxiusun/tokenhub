# override_test.go 代码阅读文档

## 1. 全局总结

本文件是 `override.go` 的单元测试，全面测试了参数覆盖系统的各种操作模式、条件执行、通配符路径、请求头操作等核心功能。

## 2. 依赖关系

- `common`: JSON 操作
- `dto`: 请求 DTO
- `types`: 错误类型

## 3. 测试用例分类

### 基础操作测试
- `TestApplyParamOverrideTrimPrefix/Suffix`: 前缀/后缀裁剪
- `TestApplyParamOverrideSet`: 设置值
- `TestApplyParamOverrideDelete`: 删除值
- `TestApplyParamOverrideMove`: 移动值
- `TestApplyParamOverrideCopy`: 复制值
- `TestApplyParamOverrideReplace/RegexReplace`: 字符串替换
- `TestApplyParamOverrideEnsurePrefix/Suffix`: 确保前缀/后缀
- `TestApplyParamOverrideTrimSpace/ToLower/ToUpper`: 字符串变换

### 通配符路径测试
- `TestApplyParamOverrideDeleteWildcardPath`: 通配符删除
- `TestApplyParamOverrideSetWildcardPath`: 通配符设置
- `TestApplyParamOverrideTrimSpaceWildcardPath`: 通配符裁剪
- `TestApplyParamOverrideSetWildcardKeepOrigin`: 通配符 + keep_origin
- `TestApplyParamOverrideTrimSpaceMultiWildcardPath`: 多层通配符

### 条件执行测试
- `TestApplyParamOverrideConditionORDefault`: OR 逻辑
- `TestApplyParamOverrideConditionAND`: AND 逻辑
- `TestApplyParamOverrideConditionInvert`: 条件取反
- `TestApplyParamOverrideConditionPassMissingKey`: 键不存在时通过
- `TestApplyParamOverrideConditionFromContext`: 从上下文读取条件
- `TestApplyParamOverrideConditionFromRequestHeaders`: 从请求头读取条件

### 请求头操作测试
- `TestApplyParamOverrideSetHeaderAndUseInLaterCondition`: 设置请求头并在后续条件中使用
- `TestApplyParamOverrideCopyHeaderFromRequestHeaders`: 从请求头复制
- `TestApplyParamOverridePassHeadersSkipsMissingHeaders`: 透传跳过缺失头
- `TestApplyParamOverrideSetHeaderMapRewritesCommaSeparatedHeader`: 映射重写逗号分隔头
- `TestApplyParamOverrideSetHeaderMapAppendsTokens`: 追加 token
- `TestApplyParamOverrideSetHeaderMapKeepOnlyDeclaredDropsUndeclaredTokens`: 仅保留声明的 token

### 其他测试
- `TestApplyParamOverrideReturnError`: 返回自定义错误
- `TestApplyParamOverridePruneObjectsByTypeString/WhereAndPath`: 对象裁剪
- `TestApplyParamOverrideNegativeIndexPath`: 负数索引
- `TestApplyParamOverrideMixedLegacyAndOperations`: 旧格式与新格式混合

## 4. 关键逻辑分析

1. **测试覆盖全面**: 覆盖了 20+ 种操作模式和各种边界条件
2. **并发安全**: 测试了并发场景下的线程安全
3. **nil 安全**: 测试了 nil 指针的安全处理

## 5. 关联文件

- `relay/common/override.go`: 被测试的源文件
