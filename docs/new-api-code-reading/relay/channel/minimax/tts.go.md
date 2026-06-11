# tts.go 代码阅读文档

## 1. 全局总结
本文件实现了 MiniMax 渠道的文本转语音（TTS）功能，包括请求/响应类型定义、TTS 响应处理和聊天补全响应处理。MiniMax 的 TTS API 支持丰富的语音配置，包括音色、语速、情感、发音字典等。

## 2. 依赖关系
- **标准库**: `encoding/hex`, `encoding/json`, `errors`, `fmt`, `io`, `net/http`, `strings`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/service` — 响应头复制
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `MiniMaxTTSRequest` 结构体
```go
type MiniMaxTTSRequest struct {
    Model             string             `json:"model"`
    Text              string             `json:"text"`
    Stream            bool               `json:"stream,omitempty"`
    StreamOptions     *StreamOptions     `json:"stream_options,omitempty"`
    VoiceSetting      VoiceSetting       `json:"voice_setting"`
    PronunciationDict *PronunciationDict `json:"pronunciation_dict,omitempty"`
    AudioSetting      *AudioSetting      `json:"audio_setting,omitempty"`
    TimbreWeights     []TimbreWeight     `json:"timbre_weights,omitempty"`
    LanguageBoost     string             `json:"language_boost,omitempty"`
    VoiceModify       *VoiceModify       `json:"voice_modify,omitempty"`
    SubtitleEnable    bool               `json:"subtitle_enable,omitempty"`
    OutputFormat      string             `json:"output_format,omitempty"`
    AigcWatermark     bool               `json:"aigc_watermark,omitempty"`
}
```
MiniMax TTS 请求结构体，支持流式输出、音色设置、发音字典、音频参数、音色混合、语言增强、字幕、水印等丰富配置。

### `StreamOptions` 结构体
```go
type StreamOptions struct {
    ExcludeAggregatedAudio bool `json:"exclude_aggregated_audio,omitempty"`
}
```
流式选项，控制是否排除聚合音频数据。

### `VoiceSetting` 结构体
```go
type VoiceSetting struct {
    VoiceID           string  `json:"voice_id"`
    Speed             float64 `json:"speed,omitempty"`
    Vol               float64 `json:"vol,omitempty"`
    Pitch             int     `json:"pitch,omitempty"`
    Emotion           string  `json:"emotion,omitempty"`
    TextNormalization bool    `json:"text_normalization,omitempty"`
    LatexRead         bool    `json:"latex_read,omitempty"`
}
```
音色设置：音色 ID、语速、音量、音调、情感、文本归一化、LaTeX 朗读开关。

### `PronunciationDict` 结构体
```go
type PronunciationDict struct {
    Tone []string `json:"tone,omitempty"`
}
```
发音字典，用于自定义音调。

### `AudioSetting` 结构体
```go
type AudioSetting struct {
    SampleRate int    `json:"sample_rate,omitempty"`
    Bitrate    int    `json:"bitrate,omitempty"`
    Format     string `json:"format,omitempty"`
    Channel    int    `json:"channel,omitempty"`
    ForceCbr   bool   `json:"force_cbr,omitempty"`
}
```
音频参数设置：采样率、比特率、格式、声道数、强制恒定比特率。

### `TimbreWeight` 结构体
```go
type TimbreWeight struct {
    VoiceID string `json:"voice_id"`
    Weight  int    `json:"weight"`
}
```
音色混合权重。

### `VoiceModify` 结构体
```go
type VoiceModify struct {
    Pitch        int    `json:"pitch,omitempty"`
    Intensity    int    `json:"intensity,omitempty"`
    Timbre       int    `json:"timbre,omitempty"`
    SoundEffects string `json:"sound_effects,omitempty"`
}
```
声音修改参数：音调、强度、音色、音效。

### `MiniMaxTTSResponse` 结构体
```go
type MiniMaxTTSResponse struct {
    Data      MiniMaxTTSData   `json:"data"`
    ExtraInfo MiniMaxExtraInfo `json:"extra_info"`
    TraceID   string           `json:"trace_id"`
    BaseResp  MiniMaxBaseResp  `json:"base_resp"`
}
```
TTS 响应结构体。

### `MiniMaxTTSData` 结构体
```go
type MiniMaxTTSData struct {
    Audio  string `json:"audio"`
    Status int    `json:"status"`
}
```
音频数据，`Audio` 字段可能是 URL 或 hex 编码的音频数据。

### `MiniMaxExtraInfo` 结构体
```go
type MiniMaxExtraInfo struct {
    UsageCharacters int64 `json:"usage_characters"`
}
```
使用量信息，基于字符数计费。

### `MiniMaxBaseResp` 结构体
```go
type MiniMaxBaseResp struct {
    StatusCode int64  `json:"status_code"`
    StatusMsg  string `json:"status_msg"`
}
```
基础响应状态。

## 4. 函数详解

### `getContentTypeByFormat(format) string`
根据音频格式返回对应的 MIME 类型：mp3→audio/mpeg, wav→audio/wav, flac→audio/flac, aac→audio/aac, pcm→audio/pcm。默认返回 `audio/mpeg`。

### `handleTTSResponse(c, resp, info) (usage, err)`
TTS 响应处理主函数：
1. 读取完整响应体
2. 解析为 `MiniMaxTTSResponse`
3. 检查 `BaseResp.StatusCode` 是否为 0（成功）
4. 检查音频数据是否为空
5. 根据音频数据类型处理：
   - URL 开头 → HTTP 302 重定向
   - 其他 → hex 解码为音频字节，直接返回 `audio/mpeg` 类型数据
6. 计算 usage：promptTokens 为估算值，totalTokens 为字符数

### `handleChatCompletionResponse(c, resp, info) (usage, err)`
聊天补全响应处理：
1. 读取完整响应体
2. 复制上游响应头到客户端
3. 直接返回原始 JSON 响应

## 5. 关键逻辑分析

### 音频数据双格式
MiniMax TTS 返回的音频数据可能是 URL（需要重定向）或 hex 编码的原始音频数据。处理器根据前缀 `"http"` 判断类型。

### Hex 编码音频
MiniMax 使用 hex 编码而非 base64 编码音频数据，处理器使用 `hex.DecodeString` 进行解码。这是 MiniMax API 的特殊设计。

### 字符计费
MiniMax TTS 使用字符数（`usage_characters`）而非 token 数进行计费，与大多数 AI 服务的计费方式不同。

## 6. 关联文件
- `adaptor.go` — 在 `ConvertAudioRequest` 中构建 TTS 请求，在 `DoResponse` 中调用 `handleTTSResponse`
- `constants.go` — 包含语音模型（`speech-*`）
