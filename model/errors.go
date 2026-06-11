package model

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDisabled      = errors.New("user is disabled")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenDisabled     = errors.New("token is disabled")
	ErrTokenExpired      = errors.New("token is expired")
	ErrTokenQuotaExceeded = errors.New("token quota exceeded")
	ErrChannelNotFound   = errors.New("channel not found")
	ErrChannelDisabled   = errors.New("channel is disabled")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrDatabase          = errors.New("database error")
	ErrNotFound          = errors.New("record not found")
)
