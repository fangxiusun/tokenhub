package dto

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []Message       `json:"messages"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	Stream           bool            `json:"stream"`
	Stop             []string        `json:"stop,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	User             string          `json:"user,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// EmbeddingRequest represents an embedding request
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// ImageGenerationRequest represents an image generation request
type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n,omitempty"`
	Size   string `json:"size,omitempty"`
}

// TokenRequest represents a token creation/update request
type TokenRequest struct {
	Name             string `json:"name"`
	ExpiredTime      int64  `json:"expired_time"`
	RemainQuota      int    `json:"remain_quota"`
	UnlimitedQuota   bool   `json:"unlimited_quota"`
	ModelLimits      string `json:"model_limits"`
	ModelLimitsEnabled bool `json:"model_limits_enabled"`
}

// ChannelRequest represents a channel creation/update request
type ChannelRequest struct {
	Name         string  `json:"name"`
	Type         int     `json:"type"`
	Key          string  `json:"key"`
	BaseURL      string  `json:"base_url"`
	Models       string  `json:"models"`
	ModelMapping string  `json:"model_mapping"`
	Groups       string  `json:"groups"`
	Tags         string  `json:"tags"`
	Weight       int     `json:"weight"`
	Status       int     `json:"status"`
	ModelRatio   float64 `json:"model_ratio"`
	ModelPrice   float64 `json:"model_price"`
}

