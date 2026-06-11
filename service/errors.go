package service

import "errors"

var (
	ErrInsufficientQuota = errors.New("insufficient quota")
	ErrNoAvailableChannel = errors.New("no available channel")
	ErrChannelDisabled   = errors.New("channel is disabled")
	ErrModelNotSupported = errors.New("model not supported")
	ErrRequestFailed     = errors.New("request failed")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrRateLimited       = errors.New("rate limited")
)
