# option.go 代码阅读文档

## 1. 全局总结

该文件实现了系统选项（Option）的读取和更新接口。支持大量配置项的验证逻辑，包括 OAuth 启用检查、计费倍率验证、面板配置验证等。

## 2. 依赖关系

- `common` — 通用工具
- `i18n` — 国际化
- `model` — Option 模型
- `setting` — 模型请求限流
- `setting/console_setting` — 控制台设置验证
- `setting/operation_setting` — 运营设置
- `setting/ratio_setting` — 倍率设置
- `setting/system_setting` — 系统设置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `OptionUpdateRequest` | 选项更新请求（key + value） |

## 4. 函数详解

### `GetOptions(c *gin.Context)`
获取所有系统选项。过滤敏感字段（Token、Secret、Key 等后缀），额外生成 `CompletionRatioMeta` 汇总。

### `UpdateOption(c *gin.Context)`
更新单个系统选项。包含大量验证逻辑：
- OAuth 启用前检查必要配置
- 倍率设置格式验证
- 控制台配置验证
- 支付合规确认检查

## 5. 关键逻辑分析

- 敏感字段过滤：以 Token、Secret、Key、secret、api_key 结尾的字段不返回
- `CompletionRatioMeta` 汇总所有计费相关选项中的模型名，生成统一的完成倍率信息
- 支付合规字段（`payment_setting.compliance_*`）禁止通过通用接口修改
- 每个 OAuth 开关启用前都会检查对应的 ClientId/ClientSecret 是否已配置
- 倍率设置（GroupRatio、ImageRatio、AudioRatio 等）有专门的验证函数

## 6. 关联文件

- `model/option.go` — Option 模型
- `setting/` — 各配置模块
