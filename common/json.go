package common

import (
	"encoding/json"
	"io"
)

// Marshal serializes the given value to JSON
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent serializes the given value to indented JSON
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal deserializes the given JSON data to the given value
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// UnmarshalJsonStr deserializes the given JSON string to the given value
func UnmarshalJsonStr(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

// DecodeJson decodes JSON from a reader to the given value
func DecodeJson(reader io.Reader, v any) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(v)
}

// GetJsonType returns the type of the JSON value
func GetJsonType(data json.RawMessage) string {
	if len(data) == 0 {
		return "null"
	}
	switch data[0] {
	case '{':
		return "object"
	case '[':
		return "array"
	case '"':
		return "string"
	case 't', 'f':
		return "boolean"
	case 'n':
		return "null"
	default:
		// Check if it's a number
		if data[0] >= '0' && data[0] <= '9' || data[0] == '-' {
			return "number"
		}
		return "unknown"
	}
}

