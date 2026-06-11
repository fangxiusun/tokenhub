# billingexpr_test.go 代码阅读文档

## 1. 全局总结
该文件是 billingexpr 包的单元测试文件，用于验证计费表达式系统的各种功能。测试覆盖了表达式编译、执行、缓存、四舍五入、结算、时间函数、图像/音频 token、上下文长度变量以及基准测试等。通过模拟不同的 token 参数和请求输入，确保计费逻辑的正确性和边界情况处理。

## 2. 依赖关系
- **math**: 数学计算，如绝对值、四舍五入等。
- **math/rand**: 生成随机数，用于模糊测试。
- **testing**: Go 标准测试框架。
- **github.com/QuantumNous/new-api/pkg/billingexpr**: 被测试的包，提供表达式编译、执行、结算等功能。

## 3. 类型定义
该文件为测试文件，未定义新的类型，但使用了 billingexpr 包中的以下类型：
- **TokenParams**: 表示 token 参数的结构体，包含 P、C、Len、CR、CC、CC1h、Img、ImgO、AI、AO 等字段。
- **BillingSnapshot**: 表示计费快照的结构体，包含计费模式、模型名、表达式字符串、哈希、组比例、估计 token 数、估计配额等字段。
- **RequestInput**: 表示请求输入的结构体，包含 Headers 和 Body 字段。
- **TraceResult**: 表示追踪结果的结构体，包含匹配的层级和成本。

## 4. 函数详解
### TestClaude_StandardTier
测试 Claude 风格表达式在标准层级（p <= 200000）下的计算结果。验证成本计算和层级匹配。

### TestClaude_LongContextTier
测试 Claude 风格表达式在长上下文层级（p > 200000）下的计算结果。

### TestClaude_BoundaryExact
测试 Claude 风格表达式在边界值（p = 200000）下的计算结果。

### TestGLM_Tier1, TestGLM_Tier2, TestGLM_Tier3
测试 GLM 风格多条件表达式在不同层级下的计算结果。

### TestSimpleExpr_NoTier
测试不使用 tier() 函数的简单表达式。

### TestMathHelpers
测试数学辅助函数（max、min 等）在表达式中的使用。

### TestRequestProbeHelpers
测试 param() 和 header() 辅助函数，用于探测请求参数和头部信息。

### TestParamProbeNestedBool, TestParamProbeArrayLength
测试嵌套 JSON 参数和数组长度探测。

### TestRequestProbeMissingFieldReturnsNil
测试探测不存在的字段时返回 nil。

### TestRequestProbeMultipleRulesMultiply
测试多个探测规则相乘的情况。

### TestCeilFloor
测试 ceil 和 floor 函数。

### TestZeroTokens
测试 token 参数为零时的计算结果。

### TestQuotaRound
测试配额四舍五入函数 QuotaRound。

### TestComputeTieredQuota_Basic, TestComputeTieredQuota_SameTier
测试分层配额计算的基本场景和相同层级场景。

### TestCompileError
测试编译错误表达式。

### TestCompileCache_SameResult, TestInvalidateCache
测试编译缓存和缓存失效。

### TestExprHashString_Deterministic
测试表达式哈希的确定性。

### TestCachePresent_StandardTier, TestCachePresent_LongContextTier
测试包含缓存 token 的表达式计算。

### TestCacheAbsent_ZeroCacheTokens
测试缓存 token 为零时的表达式计算。

### TestMixedCacheFields, TestMixedCacheFields_AllCacheZero
测试混合缓存字段和全部缓存为零的情况。

### TestBackwardCompat_OldExprWithTokenParams
测试旧表达式与新 TokenParams 的向后兼容性。

### TestComputeTieredQuota_WithCache, TestComputeTieredQuota_WithCacheCrossTier
测试带缓存 token 的分层配额计算和跨层级场景。

### TestFuzz_NonNegativeResults, TestFuzz_SettlementConsistency
模糊测试，验证结果非负和结算一致性。

### TestComputeTieredQuota_BasicSettlement, TestComputeTieredQuota_WithGroupRatio, TestComputeTieredQuota_ZeroTokens, TestComputeTieredQuota_RoundingEdge, TestComputeTieredQuota_RoundingEdgeDown
测试分层配额计算的各种边界情况。

### TestComputeTieredQuotaWithRequest_ProbeAffectsQuota, TestComputeTieredQuota_BoundaryTierCrossing
测试探测影响配额和边界层级跨越。

### TestTimeFunctions_*
测试时间函数（hour、minute、weekday、month、day）在不同时区下的行为。

### TestImageTokenVariable, TestAudioTokenVariables, TestImageAudioVariables, TestImageAudioZero
测试图像和音频 token 变量。

### TestLen_*, TestLen_BoundaryExact, TestLen_BoundaryPlusOne, TestLen_ZeroDefaultsToZero
测试上下文长度变量 len 在不同层级下的行为。

### BenchmarkExprCompile, BenchmarkExprRunCached
基准测试，比较编译和缓存执行的性能。

## 5. 关键逻辑分析
- **表达式语言**: 测试覆盖了 Claude 和 GLM 风格的计费表达式，包括层级条件、数学函数、时间函数、参数探测等。
- **缓存机制**: 通过重复执行相同表达式验证缓存命中，并测试缓存失效。
- **结算逻辑**: 测试分层配额计算，包括组比例、四舍五入、层级跨越检测。
- **模糊测试**: 使用随机 token 参数验证结果非负和结算一致性，增强鲁棒性。
- **边界条件**: 重点测试边界值（如 p=200000）、零 token、缺失字段等场景。

## 6. 关联文件
- **compile.go**: 提供表达式编译和缓存功能，测试中通过 RunExpr 间接调用。
- **run.go**: 提供表达式执行功能，测试中直接调用 RunExpr 和 RunExprWithRequest。
- **settle.go**: 提供分层配额计算功能，测试中调用 ComputeTieredQuota 和 ComputeTieredQuotaWithRequest。
- **types.go**: 定义测试中使用的类型（TokenParams、BillingSnapshot 等）。
- **round.go**: 提供配额四舍五入函数，测试中调用 QuotaRound。