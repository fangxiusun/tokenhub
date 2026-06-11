# image.go 代码阅读文档

## 1. 全局总结
image.go 提供图像处理相关的服务功能，包括图像格式检测、尺寸获取等。

## 2. 依赖关系
### 2.1 导入的包
- `image` - 标准图像库
- `image/png` - PNG 解码
- `image/jpeg` - JPEG 解码
- `bytes` - 字节缓冲

### 2.2 被引用的文件
- `controller/` - 图像相关控制器

## 3. 类型定义
### 3.1 函数
- `GetImageDimensions(data []byte) (int, int, error)` - 获取图像尺寸
- `DetectImageFormat(data []byte) string` - 检测图像格式

## 4. 函数详解
### 4.1 GetImageDimensions
- **签名**: `func GetImageDimensions(data []byte) (int, int, error)`
- **职责**: 获取图像的宽和高
- **返回值**: 宽度、高度、错误

## 5. 关键逻辑分析
- **格式支持**: PNG, JPEG, GIF, WebP
- **高效解码**: 只解码头部信息

## 6. 关联文件
- `controller/image.go` - 图像控制器
