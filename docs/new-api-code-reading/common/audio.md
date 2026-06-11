# audio.go 代码阅读文档

## 1. 全局总结

`audio.go` 是一个纯 Go 实现的音频文件时长解析工具集，提供获取多种常见音频格式文件时长的功能。该文件完全不依赖外部的 `ffmpeg` 或 `ffprobe` 程序，而是通过多个第三方 Go 音频解析库来解析不同格式的音频文件。支持的格式包括 MP3、WAV、FLAC、M4A/MP4、OGG/Vorbis、Opus、AIFF、WebM、AAC 共 9 种格式。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/abema/go-mp4` | 解析 M4A/MP4 容器格式 |
| `github.com/go-audio/aiff` | 解析 AIFF 音频格式 |
| `github.com/go-audio/wav` | 解析 WAV 音频格式 |
| `github.com/jfreymuth/oggvorbis` | 解析 OGG/Vorbis 音频格式 |
| `github.com/mewkiz/flac` | 解析 FLAC 音频格式 |
| `github.com/pkg/errors` | 错误包装，提供上下文信息 |
| `github.com/tcolgate/mp3` | 解析 MP3 音频格式 |
| `github.com/yapingcat/gomedia/go-codec` | 解析 AAC ADTS 帧 |
| `context` | 标准库，支持上下文取消 |
| `encoding/binary` | 标准库，二进制数据解析（OGG granule position） |
| `fmt` | 标准库，格式化输出 |
| `io` | 标准库，I/O 接口 |

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### 主入口函数

#### `GetAudioDuration(ctx context.Context, f io.ReadSeeker, ext string) (float64, error)`

- **功能**：根据文件扩展名分发到对应的格式解析函数
- **参数**：
  - `ctx context.Context` — 上下文，预留用于超时控制（当前未实际使用）
  - `f io.ReadSeeker` — 音频文件的读取器，需要支持 Seek 操作
  - `ext string` — 文件扩展名（如 `.mp3`、`.wav` 等）
- **返回值**：`(duration float64, err error)` — 音频时长（秒）和错误
- **支持的格式**：
  - `.mp3`
  - `.wav`
  - `.flac`
  - `.m4a`、`.mp4`
  - `.ogg`、`.oga`、`.opus`（OGG 失败时尝试 Opus 解析）
  - `.aiff`、`.aif`、`.aifc`
  - `.webm`
  - `.aac`

### 私有解析函数

#### `getMP3Duration(r io.Reader) (float64, error)`

- **原理**：逐帧解码 MP3 文件，累加每帧时长
- **注意**：VBR（可变比特率）MP3 的估算可能不完全精确，但通常足够好

#### `getWAVDuration(r io.ReadSeeker) (float64, error)`

- **原理**：解析 WAV 文件头获取采样率、声道数、位深，然后通过 PCM 数据大小计算时长
- **特殊处理**：
  - 强制复位文件指针到开头
  - 当 `PCMSize` 为 0 时，通过文件大小反推数据区大小
  - 使用 `bytesPerFrame = numChans * (bitDepth / 8)` 计算帧大小

#### `getFLACDuration(r io.Reader) (float64, error)`

- **原理**：解析 FLAC STREAMINFO 块，获取总采样数和采样率，`时长 = 总采样数 / 采样率`

#### `getM4ADuration(r io.ReadSeeker) (float64, error)`

- **原理**：使用 `mp4.Probe` 探测文件，获取 Duration 和 Timescale
- **计算公式**：`时长 = Duration / Timescale`

#### `getOGGDuration(r io.ReadSeeker) (float64, error)`

- **原理**：创建 OGG Vorbis Reader，读取全部采样数据，通过总采样数 / 采样率计算时长
- **注意**：需要读取整个文件来获取总采样数

#### `getOpusDuration(r io.ReadSeeker) (float64, error)`

- **原理**：手动解析 OGG 页面头部，提取 granule position（累计采样数）
- **特殊处理**：
  - Opus 固定采样率为 48000 Hz
  - 通过检查 "OggS" 标识找到 OGG 页面
  - 使用小端序读取 granule position（字节 6-13）

#### `getAIFFDuration(r io.ReadSeeker) (float64, error)`

- **原理**：使用 `aiff.NewDecoder` 解析 AIFF 文件，调用 `Duration()` 方法获取时长

#### `getWebMDuration(r io.ReadSeeker) (float64, error)`

- **原理**：简化的实现，检查 EBML 标识但**不支持**完整解析
- **当前状态**：始终返回错误，提示需要完整的 EBML 解析器或使用 ffprobe

#### `getAACDuration(r io.ReadSeeker) (float64, error)`

- **原理**：读取整个文件，使用 `gomedia` 的 `SplitAACFrame` 分割帧，解析 ADTS 头部获取采样率
- **计算公式**：`时长 = (帧数 × 1024) / 采样率`（每个 AAC ADTS 帧包含 1024 个采样）

## 5. 关键逻辑分析

1. **纯 Go 实现**：该文件完全使用 Go 库解析音频，无需外部二进制依赖，简化了部署流程。

2. **格式降级策略**：OGG 格式解析失败时，会尝试 Opus 解析（第35-37行），因为 `.ogg`/`.oga`/`.opus` 共享同一扩展名空间。

3. **WebM 不完整支持**：`getWebMDuration` 是一个占位实现，仅检查 EBML 标识但不实际解析 Duration，说明 WebM 格式支持尚不完善。

4. **错误包装**：所有错误都通过 `errors.Wrap` 包装，提供调用链上下文，便于调试。

5. **Seek 操作**：多个解析函数开头都有 `r.Seek(0, io.SeekStart)` 复位操作，确保从文件头开始解析，避免前一个解析器残留的读取位置影响。

6. **资源泄漏风险**：`getMP3Duration` 中的 `mp3.Frame` 没有显式关闭，但该库通常不需要手动关闭。`getFLACDuration` 中的 `stream.Close()` 使用 `defer` 确保资源释放。

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `common/sys_log.go` | `SysLog` 函数用于日志记录 |
| `relay/utils/audio.go` 或类似文件 | 可能在中继层调用此函数解析音频时长 |
| `controller/upload.go` | 上传处理时可能调用此函数获取音频时长 |
