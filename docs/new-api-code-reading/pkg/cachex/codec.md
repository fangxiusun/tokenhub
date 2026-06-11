# codec.go 代码阅读文档

## 1. 全局总结
该文件定义了缓存值编解码器接口和实现，用于在缓存中存储和检索不同类型的数据。提供 IntCodec、StringCodec 和 JSONCodec 三种实现，分别处理整数、字符串和任意类型的 JSON 序列化。

## 2. 依赖关系
- **encoding/json**: JSON 序列化和反序列化。
- **fmt**: 格式化错误信息。
- **strconv**: 整数与字符串转换。
- **strings**: 字符串处理（去除空格）。

## 3. 类型定义
### ValueCodec[V any]
泛型接口，定义编解码方法：
- **Encode(v V) (string, error)**: 将值编码为字符串。
- **Decode(s string) (V, error)**: 将字符串解码为值。

### IntCodec
结构体，实现 ValueCodec[int]，用于整数编解码。

### StringCodec
结构体，实现 ValueCodec[string]，用于字符串编解码。

### JSONCodec[V any]
泛型结构体，实现 ValueCodec[V]，用于任意类型的 JSON 编解码。

## 4. 函数详解
### IntCodec.Encode
将整数转换为字符串。

### IntCodec.Decode
将字符串转换为整数，去除空格，空字符串返回错误。

### StringCodec.Encode
直接返回输入字符串。

### StringCodec.Decode
直接返回输入字符串。

### JSONCodec.Encode
使用 json.Marshal 将值序列化为 JSON 字符串。

### JSONCodec.Decode
使用 json.Unmarshal 将 JSON 字符串反序列化为空白字符串返回错误。

## 5. 关键逻辑分析
- **泛型设计**: 使用 Go 泛型实现类型安全的编解码器。
- **错误处理**: 空值或无效格式返回明确错误。
- **灵活性**: JSONCodec 支持任意可序列化类型，适用于复杂缓存值。

## 6. 关联文件
- **hybrid_cache.go**: 使用 ValueCodec 接口进行 Redis 编解码。
- **namespace.go**: 提供键命名空间，与编解码器配合使用。