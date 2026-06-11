# 配置管理层详细设计 (`setting/`)

## 1. 概述
配置管理层包含所有系统功能的配置管理，组织成子包。每个子包负责一个特定功能领域的配置。

## 2. 文件详细说明

### 2.1 根目录文件
- **`auto_group.go`** -- 自动组分配管理：管理自动分配的用户组列表（如 `default`），带 JSON 序列化。
- **`chat.go`** -- 预配置的聊天客户端深度链接：Cherry Studio, Lobe Chat, OpenCat, NextChat 等的 URL 模板。
- **`midjourney.go`** -- Midjourney 功能开关标志：通知、账户过滤、模式清除、转发 URL、操作检查。
- **`payment_creem.go`** -- Creem 支付网关设置：API 密钥、产品、测试模式、webhook 密钥。
- **`payment_stripe.go`** -- Stripe 支付网关设置：API 密钥、webhook 密钥、价格 ID、单价、最低充值、促销代码。
- **`payment_waffo.go`** -- Waffo 支付网关设置：API 密钥、证书、商户 ID、货币、定价、沙箱模式。
- **`payment_waffo_pancake.go`** -- Waffo Pancake 托管结账配置：商户 ID、私钥、商店/产品 ID、定价。
- **`rate_limit.go`** -- 模型级请求限流：每模型限制组、持续时间和计数设置，带 JSON 同步。
- **`sensitive.go`** -- 敏感词检测设置：启用/禁用标志、敏感时停止开关、流缓存队列长度和敏感词列表。
- **`user_usable_group.go`** -- 用户可用组管理：组名到显示名的线程安全映射，带 JSON 序列化。

### 2.2 `billing_setting/` 子包
- **`tiered_billing.go`** -- 分层/基于表达式的计费系统：每模型计费模式（`ratio` 或 `tiered_expr`）和表达式配置。从 `pkg/billingexpr` 解析计费表达式。

### 2.3 `config/` 子包
- **`config.go`** -- `ConfigManager` -- 统一配置注册表：命名配置模块的线程安全映射，基于 JSON 的加载/保存，类型安全的 getter（`Get[T]`），以及 DB 同步集成。所有设置子包通过 `GlobalConfig.Register()` 注册。

### 2.4 `console_setting/` 子包
- **`config.go`** -- 控制台/仪表盘设置：API 信息面板、Uptime Kuma 组、公告、FAQ -- 每个都有启用/禁用开关。
- **`validation.go`** -- 控制台设置的验证逻辑：URL 验证、XSS 防护、颜色验证、slug 验证，以及公告/FAQ/API 信息的 JSON 数组解析。

### 2.5 `model_setting/` 子包
- **`claude.go`** -- Claude 特定模型设置：每模型头注入、默认最大令牌、思考适配器启用/预算百分比。
- **`gemini.go`** -- Gemini 特定模型设置：安全设置、版本设置、imagine 模型、思考适配器、函数调用思考签名。
- **`global.go`** -- 全局模型设置：`ChatCompletionsToResponsesPolicy`（为特定渠道/模型自动将聊天补全转发到 Responses API），带渠道/类型/模式匹配。
- **`grok.go`** -- Grok 模型设置：违规扣减启用/金额。
- **`qwen.go`** -- Qwen 模型设置：Qwen/阿里云提供商的图像同步模型列表。

### 2.6 `operation_setting/` 子包
- **`operation_setting.go`** -- 核心运营设置：演示站点开关、自使用模式、上游错误自动禁用关键词。
- **`channel_affinity_setting.go`** -- 渠道亲和性设置（粘性路由）。
- **`checkin_setting.go`** -- 每日签到奖励设置。
- **`general_setting.go`** -- 通用运营设置（站点名称、页脚、Logo 等）。
- **`monitor_setting.go`** -- 渠道监控/告警设置。
- **`payment_setting.go`** -- 支付运营设置（定价、充值等）。
- **`payment_setting_old.go`** -- 旧版支付设置（向后兼容）。
- **`quota_setting.go`** -- 配额管理设置（预扣、配额单位等）。
- **`status_code_ranges.go`** -- HTTP 状态码范围定义用于重试/错误分类。
- **`token_setting.go`** -- 令牌生成/显示设置（默认令牌生成、令牌长度等）。
- **`tools.go`** -- 运营工具实用程序。

### 2.7 `perf_metrics_setting/` 子包
- **`config.go`** -- 性能指标设置：启用/禁用、刷新间隔、桶时间（分钟/小时/天）、指标数据保留天数。

### 2.8 `performance_setting/` 子包
- **`config.go`** -- 性能调优设置：磁盘缓存（启用、阈值、最大大小、路径）和系统监控阈值（CPU、内存、磁盘）。同步到 `common.DiskCacheConfig` 和 `common.PerformanceMonitorConfig`。

### 2.9 `ratio_setting/` 子包
- **`cache_ratio.go`** -- 每模型默认缓存命中折扣比例（如 GPT-4o 为 0.5x，Gemini 为 0.1x）。
- **`compact_suffix.go`** -- 紧凑模型后缀工具：为 OpenAI 紧凑模式响应附加 `-openai-compact`。
- **`expose_ratio.go`** -- 向非管理员用户公开模型比例/定价数据的开关。
- **`exposed_cache.go`** -- 暴露比例数据的基于 TTL 的缓存。
- **`group_ratio.go`** -- 组比例管理：每组定价乘数、组到组比例和特殊可用组映射。
- **`model_ratio.go`** -- **核心定价表**：750+ 行的默认模型到比例映射（所有 OpenAI, Claude, Gemini, Qwen, DeepSeek 等模型），基准单位 `$0.002 / 1K tokens = ratio 1`。

### 2.10 `reasoning/` 子包
- **`suffix.go`** -- 推理努力后缀解析：从模型名称中剥离 `-high`/`-low`/`-medium`/`-max`/`-minimal`/`-none` 后缀，并将它们映射到 OpenAI 和 DeepSeek 的推理努力级别。

### 2.11 `system_setting/` 子包
- **`discord.go`** -- Discord OAuth 设置：客户端 ID 和密钥。
- **`fetch_setting.go`** -- SSRF 防护设置：启用标志、域名/IP 过滤模式和列表、允许端口。
- **`legal.go`** -- 法律页面内容：用户协议和隐私政策文本。
- **`oidc.go`** -- OpenID Connect (OIDC) 设置：客户端 ID/密钥、well-known URL、授权/令牌/用户信息端点。
- **`passkey.go`** -- WebAuthn/Passkey 设置：RP 显示名称、RP ID、来源、用户验证、附件偏好。
- **`theme.go`** -- 前端主题选择：同步主题名称到 `common.SetTheme()`。
- **`system_setting_old.go`** -- 旧版系统设置（向后兼容）。

---

## 关联文件列表

### 配置管理层核心文件
- `setting/auto_group.go` - 自动组管理
- `setting/chat.go` - 聊天客户端链接
- `setting/midjourney.go` - Midjourney 设置
- `setting/payment_creem.go` - Creem 支付设置
- `setting/payment_stripe.go` - Stripe 支付设置
- `setting/payment_waffo.go` - Waffo 支付设置
- `setting/payment_waffo_pancake.go` - Waffo-Pancake 设置
- `setting/rate_limit.go` - 模型限流设置
- `setting/sensitive.go` - 敏感词设置
- `setting/user_usable_group.go` - 用户组设置

### 子包文件
- `setting/billing_setting/tiered_billing.go` - 分层计费设置
- `setting/config/config.go` - 配置管理器
- `setting/console_setting/config.go` - 控制台设置
- `setting/console_setting/validation.go` - 控制台设置验证
- `setting/model_setting/claude.go` - Claude 模型设置
- `setting/model_setting/gemini.go` - Gemini 模型设置
- `setting/model_setting/global.go` - 全局模型设置
- `setting/model_setting/grok.go` - Grok 模型设置
- `setting/model_setting/qwen.go` - Qwen 模型设置
- `setting/operation_setting/operation_setting.go` - 核心运营设置
- `setting/operation_setting/channel_affinity_setting.go` - 渠道亲和性设置
- `setting/operation_setting/checkin_setting.go` - 签到设置
- `setting/operation_setting/general_setting.go` - 通用运营设置
- `setting/operation_setting/monitor_setting.go` - 监控设置
- `setting/operation_setting/payment_setting.go` - 支付设置
- `setting/operation_setting/payment_setting_old.go` - 旧版支付设置
- `setting/operation_setting/quota_setting.go` - 配额设置
- `setting/operation_setting/status_code_ranges.go` - 状态码范围
- `setting/operation_setting/token_setting.go` - 令牌设置
- `setting/operation_setting/tools.go` - 运营工具
- `setting/perf_metrics_setting/config.go` - 性能指标设置
- `setting/performance_setting/config.go` - 性能调优设置
- `setting/ratio_setting/cache_ratio.go` - 缓存比例
- `setting/ratio_setting/compact_suffix.go` - 紧凑后缀
- `setting/ratio_setting/expose_ratio.go` - 比例公开
- `setting/ratio_setting/exposed_cache.go` - 暴露缓存
- `setting/ratio_setting/group_ratio.go` - 组比例
- `setting/ratio_setting/model_ratio.go` - 模型比例
- `setting/reasoning/suffix.go` - 推理后缀
- `setting/system_setting/discord.go` - Discord 设置
- `setting/system_setting/fetch_setting.go` - SSRF 设置
- `setting/system_setting/legal.go` - 法律页面
- `setting/system_setting/oidc.go` - OIDC 设置
- `setting/system_setting/passkey.go` - Passkey 设置
- `setting/system_setting/theme.go` - 主题设置
- `setting/system_setting/system_setting_old.go` - 旧版系统设置

### 依赖的包
- `pkg/billingexpr/` - 计费表达式引擎
- `common/` - 公共工具
