# constant.go 代码阅读文档

## 1. 全局总结
本文件定义了 OpenRouter 渠道的常量，包括渠道名称和模型列表。作为渠道的基础配置文件，为路由系统提供渠道标识和可用模型清单。

## 2. 依赖关系
无外部依赖，仅使用 Go 内置类型。

## 3. 类型定义

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | OpenRouter 支持的模型列表，当前为空切片（待填充） |
| `ChannelName` | `string` | 渠道名称常量，值为 `"openrouter"` |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- `ModelList` 当前为空切片，表明该渠道的模型列表可能在其他地方动态填充，或尚未完全配置。
- `ChannelName` 用于在系统中唯一标识 OpenRouter 渠道。

## 6. 关联文件
- `relay/channel/openrouter/dto.go` — 同包下的数据传输对象定义
- `relay/channel/adapter.go` — `Adaptor` 接口定义，所有渠道适配器必须实现该接口
