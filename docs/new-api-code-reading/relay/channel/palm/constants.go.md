# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 Google PaLM 渠道的常量配置，包括支持的模型列表和渠道名称。为路由系统提供渠道标识信息。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | PaLM 支持的模型列表，包含 `"PaLM-2"` |
| `ChannelName` | `string` | 渠道名称，值为 `"google palm"` |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- 该渠道仅支持单一模型 `PaLM-2`，与 `adaptor.go` 中固定的 `chat-bison-001` 端点对应。
- 渠道名称 `"google palm"` 用于在管理界面和日志中标识该渠道。

## 6. 关联文件
- `relay/channel/palm/adaptor.go` — 适配器实现，使用这些常量
- `relay/channel/palm/dto.go` — PaLM API 数据结构定义
