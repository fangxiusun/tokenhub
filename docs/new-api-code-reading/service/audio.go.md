# audio.go 代码阅读文档

## 1. 全局总结
audio.go 提供音频处理相关的服务功能，包括音频时长提取、格式检测等。

## 2. 依赖关系
### 2.1 导入的包
- `common` - 公共工具

### 2.2 被引用的文件
- `controller/` - 音频相关控制器

## 3. 类型定义
### 3.1 函数
- `GetAudioDuration(data []byte) (float64, error)` - 获取音频时长
- `DetectAudioFormat(data []byte) string` - 检测音频格式

## 4. 函数详解
### 4.1 GetAudioDuration
- **签名**: `func GetAudioDuration(data []byte) (float64, error)`
- **职责**: 从音频数据中提取时长信息
- **支持格式**: MP3, WAV, FLAC, OGG, M4A

### 4.2 DetectAudioFormat
- **签名**: `func DetectAudioFormat(data []byte) string`
- **职责**: 根据文件头检测音频格式

## 5. 关键逻辑分析
- **多格式支持**: 通过魔数（magic bytes）检测格式
- **纯 Go 实现**: 无需外部依赖

## 6. 关联文件
- `common/audio.go` - 音频工具函数
