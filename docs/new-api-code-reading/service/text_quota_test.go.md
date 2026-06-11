# text_quota_test.go 代码阅读文档

## 1. 全局总结
text_quota_test.go 测试文本配额计算功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `service/text_quota.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestCalculateQuota` - 测试配额计算
- `TestModelRatio` - 测试模型比例

## 4. 函数详解
### 4.1 TestCalculateQuota
- **职责**: 测试基于 Token 数的配额计算
- **验证点**: 输入/输出 Token、缓存折扣

### 4.2 TestModelRatio
- **职责**: 测试不同模型的计费比例
- **验证点**: 默认比例、自定义比例

## 5. 关键逻辑分析
- **计费公式**: quota = (input + output * completionRatio) * modelRatio * groupRatio

## 6. 关联文件
- `service/text_quota.go` - 被测试的文件
