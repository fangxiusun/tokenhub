# validate.go 代码阅读文档

## 1. 全局总结
validate.go 初始化全局结构体验证器，基于 go-playground/validator 库。

## 2. 依赖关系
### 2.1 导入的包
- `github.com/go-playground/validator/v10` - 结构体验证库

### 2.2 被引用的文件
- `common/` - 公共工具
- `controller/` - 控制器验证请求

## 3. 类型定义
### 3.1 变量
- `Validate *validator.Validate` - 全局验证器实例

## 4. 函数详解
### 4.1 init
- **职责**: 初始化全局验证器
- **逻辑流程**: 创建 validator.Validate 实例并赋值给全局变量

## 5. 关键逻辑分析
- **全局单例**: 使用包级变量实现单例模式
- **自定义验证**: 可注册自定义验证规则

## 6. 关联文件
- `controller/` - 使用 Validate 进行请求验证
