# vendor_meta.go 代码阅读文档

## 1. 全局总结

该文件实现了供应商（Vendor）元数据的 CRUD 管理接口。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 供应商模型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

- `GetAllVendors` — 分页获取所有供应商
- `SearchVendors` — 搜索供应商
- `GetVendor` — 获取单个供应商
- `CreateVendor` — 创建供应商
- `UpdateVendor` — 更新供应商
- `DeleteVendor` — 删除供应商

## 5. 关键逻辑分析

- 名称不能为空
- 创建/更新时检查名称唯一性

## 6. 关联文件

- `model/vendor.go` — 供应商模型
