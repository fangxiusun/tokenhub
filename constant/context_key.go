package constant

// Context keys
const (
	ContextKeyUserId      = "userId"
	ContextKeyUsername    = "username"
	ContextKeyRole        = "role"
	ContextKeyTokenId     = "tokenId"
	ContextKeyTokenName   = "tokenName"
	ContextKeyTokenQuota  = "tokenQuota"
	ContextKeyUserGroup   = "userGroup"
	ContextKeyChannelId   = "channelId"
	ContextKeyChannelName = "channelName"
	ContextKeyChannelType = "channelType"
	ContextKeyModel       = "model"
	ContextKeyRequestId   = "requestId"
	ContextKeyRequestBody = "requestBody"
	ContextKeyIsStream    = "isStream"
	ContextKeyRelayMode   = "relayMode"
	ContextKeyPromptTokens     = "promptTokens"
	ContextKeyCompletionTokens = "completionTokens"
	ContextKeyStartTime        = "startTime"
	ContextKeyEndTime          = "endTime"
)

// Roles
const (
	RoleGuest  = 0
	RoleUser   = 1
	RoleAdmin  = 10
	RoleRoot   = 100
)
