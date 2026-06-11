# task.go 代码阅读文档

## 1. 全局总结

该文件提供任务平台和动作到模型名称的转换功能。是一个非常简洁的工具文件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `constant` | 任务平台常量 |

## 3. 函数详解

### `CoverTaskActionToModelName(platform, action) string`
将平台和动作组合为模型名称（如 `suno_generate`）

## 4. 关联文件

- `constant/task.go` — 任务平台常量
- `task_billing.go` — 任务计费
