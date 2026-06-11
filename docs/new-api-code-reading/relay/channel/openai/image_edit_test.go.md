# image_edit_test.go 代码阅读文档

## 1. 全局总结
image_edit_test.go 测试 OpenAI 图像编辑功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `relay/channel/openai/relay_image.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestImageEdit` - 测试图像编辑请求
- `TestImageVariation` - 测试图像变体请求

## 4. 函数详解
### 4.1 TestImageEdit
- **职责**: 测试图像编辑 API 调用
- **验证点**: 请求格式、参数传递

## 5. 关键逻辑分析
- **多模态**: 测试图像输入/输出处理

## 6. 关联文件
- `relay/channel/openai/relay_image.go` - 被测试的文件
