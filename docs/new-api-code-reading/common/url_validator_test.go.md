# url_validator_test.go 代码阅读文档

## 1. 全局总结
url_validator_test.go 是 URL 验证器的测试文件，测试重定向 URL 验证功能。

## 2. 依赖关系
### 2.1 导入的包
- `testing` - 测试框架

### 2.2 被引用的文件
- `common/url_validator.go` - 被测试的文件

## 3. 类型定义
### 3.1 测试函数
- `TestValidateRedirectURL` - 测试重定向 URL 验证
- `TestIsTrustedDomain` - 测试受信任域名检查

## 4. 函数详解
### 4.1 TestValidateRedirectURL
- **职责**: 测试各种 URL 的验证结果
- **测试用例**:
  - http/https URL 应通过
  - ftp URL 应失败
  - 受信任域名应通过
  - 非受信任域名应失败

### 4.2 TestIsTrustedDomain
- **职责**: 测试域名匹配逻辑
- **测试用例**:
  - 精确匹配
  - 子域名匹配
  - 通配符匹配

## 5. 关键逻辑分析
- **边界测试**: 覆盖各种边界情况
- **安全测试**: 测试安全相关的验证逻辑

## 6. 关联文件
- `common/url_validator.go` - 被测试的文件
