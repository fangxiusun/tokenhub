# subscription_payment_epay.go 代码阅读文档

## 1. 全局总结
subscription_payment_epay.go 实现了订阅的 Epay 支付处理。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具
- `model` - 数据模型
- `service` - 业务逻辑
- `gin-gonic/gin` - Web 框架

### 2.2 被引用的文件
- `router/api-router.go` - 管理 API 路由

## 3. 类型定义
### 3.1 函数
- `SubscriptionRequestEpay` - 请求 Epay 支付
- `SubscriptionEpayNotify` - Epay 支付回调

## 4. 函数详解
### 4.1 SubscriptionRequestEpay
- **职责**: 创建 Epay 支付订单
- **逻辑流程**:
  1. 验证用户订阅计划
  2. 创建支付订单
  3. 生成支付链接
  4. 返回支付链接

### 4.2 SubscriptionEpayNotify
- **职责**: 处理 Epay 支付回调
- **逻辑流程**:
  1. 验证签名
  2. 更新订单状态
  3. 激活订阅

## 5. 关键逻辑分析
- **签名验证**: 防止伪造回调
- **幂等处理**: 重复回调不重复处理

## 6. 关联文件
- `model/subscription.go` - 订阅模型
- `service/epay.go` - Epay 支付服务
