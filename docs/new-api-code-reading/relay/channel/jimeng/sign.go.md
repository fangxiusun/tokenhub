# sign.go 代码阅读文档

## 1. 全局总结
本文件实现了即梦（Jimeng）API 的 HMAC-SHA256 签名认证机制。即梦使用类似 AWS SigV4 的签名方案，对请求进行签名以确保安全认证。签名过程包括构建规范化请求、计算 payload hash、构造签名字符串、多级 HMAC 派生签名密钥。

## 2. 依赖关系
- **标准库**: `bytes`, `crypto/hmac`, `crypto/sha256`, `encoding/hex`, `encoding/json`, `errors`, `fmt`, `io`, `net/http`, `net/url`, `sort`, `strings`, `time`
- **项目内部**:
  - `github.com/QuantumNous/new-api/logger` — 日志
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### 常量
- `HexPayloadHashKey = "HexPayloadHash"` — Gin context 中存储 payload hash 的键名

## 4. 函数详解

### `SetPayloadHash(c, req) error`
预计算请求体的 SHA256 哈希值并存入 Gin context：
1. 将请求体序列化为 JSON
2. 计算 SHA256 哈希
3. 将十六进制哈希存入 context 的 `HexPayloadHashKey` 键

### `getPayloadHash(c) string`
从 Gin context 中获取预计算的 payload hash。

### `Sign(c, req, apiKey) error`
核心签名函数，实现 HMAC-SHA256 签名流程：
1. 读取请求体并计算 SHA256 哈希（payload hash）
2. 解析 API Key（格式：`accessKey|secretKey`）
3. 生成时间戳：`X-Date`（`20060102T150405Z`）和 `shortDate`（`20060102`）
4. 构建规范化查询字符串（排序参数）
5. 构建规范化请求头（host, x-date, x-content-sha256, content-type）
6. 构建规范化请求字符串：`METHOD\nPATH\nQUERY\nHEADERS\nSIGNED_HEADERS\nPAYLOAD_HASH`
7. 计算规范化请求的 SHA256 哈希
8. 构造签名字符串：`HMAC-SHA256\nDATE\nCREDENTIAL_SCOPE\nHASHED_CANONICAL_REQUEST`
9. 多级 HMAC 派生签名密钥：`secretKey → date → region → service → signing`
10. 计算最终签名
11. 设置 `Authorization` 头：`HMAC-SHA256 Credential=..., SignedHeaders=..., Signature=...`

### `hmacSHA256(key, data) []byte`
计算 HMAC-SHA256。

## 5. 关键逻辑分析

### AWS SigV4 风格签名
即梦的签名方案与 AWS SigV4 非常相似：
- 使用 `HMAC-SHA256` 算法
- 多级密钥派生：`date → region → service → request`
- 固定区域：`cn-north-1`
- 固定服务名：`cv`（计算机视觉）

### API Key 格式
API Key 采用 `ak|sk` 格式，用 `|` 分隔 accessKey 和 secretKey。如果格式不正确会返回错误。

### 规范化请求
签名前对请求进行严格规范化：
- 查询参数按键名排序
- 请求头小写化并排序
- payload 使用 SHA256 哈希而非原始内容

### 签名范围
签名覆盖 4 个头字段：`host`, `x-date`, `x-content-sha256`, `content-type`。如果 Content-Type 为空，自动设置为 `application/json`。

## 6. 关联文件
- `adaptor.go` — 在 `DoRequest` 中调用 `Sign` 函数
- `image.go` — 依赖签名后的请求进行图像生成
