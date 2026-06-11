package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// OAuthBinding stores OAuth provider bindings for users
type OAuthBinding struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	UserId    int       `json:"user_id" gorm:"index"`
	Provider  string    `json:"provider" gorm:"size:32;index"`
	OpenId    string    `json:"open_id" gorm:"size:128"`
	CreatedAt time.Time `json:"created_at"`
}

// Supported OAuth providers
var SupportedOAuthProviders = []string{
	"github",
	"discord",
	"oidc",
	"telegram",
	"wechat",
}

// GetOAuthBindingByUserIdAndProvider returns OAuth binding by user ID and provider
func GetOAuthBindingByUserIdAndProvider(userId int, provider string) (*OAuthBinding, error) {
	var binding OAuthBinding
	err := DB.Where("user_id = ? AND provider = ?", userId, provider).First(&binding).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &binding, nil
}

// GetOAuthBindingByProviderAndOpenId returns OAuth binding by provider and open ID
func GetOAuthBindingByProviderAndOpenId(provider, openId string) (*OAuthBinding, error) {
	var binding OAuthBinding
	err := DB.Where("provider = ? AND open_id = ?", provider, openId).First(&binding).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &binding, nil
}

// CreateOAuthBinding creates a new OAuth binding
func CreateOAuthBinding(binding *OAuthBinding) error {
	return DB.Create(binding).Error
}

// DeleteOAuthBinding deletes an OAuth binding
func DeleteOAuthBinding(userId int, provider string) error {
	return DB.Where("user_id = ? AND provider = ?", userId, provider).Delete(&OAuthBinding{}).Error
}

// GetUserOAuthBindings returns all OAuth bindings for a user
func GetUserOAuthBindings(userId int) ([]OAuthBinding, error) {
	var bindings []OAuthBinding
	err := DB.Where("user_id = ?", userId).Find(&bindings).Error
	return bindings, err
}

// GetUserIdByOAuth returns user ID by OAuth provider and open ID
func GetUserIdByOAuth(provider, openId string) (int, error) {
	var binding OAuthBinding
	err := DB.Where("provider = ? AND open_id = ?", provider, openId).First(&binding).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return binding.UserId, nil
}

