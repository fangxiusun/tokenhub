package model

import (
	"time"
)

// Token represents an API token
type Token struct {
	Id            int       `json:"id" gorm:"primaryKey"`
	UserId        int       `json:"user_id" gorm:"index"`
	Key           string    `json:"key" gorm:"uniqueIndex;size:64"`
	Status        int       `json:"status" gorm:"default:1"`
	Name          string    `json:"name" gorm:"size:64"`
	ExpiredTime   int64     `json:"expired_time" gorm:"default:-1"`
	RemainQuota   int       `json:"remain_quota" gorm:"default:0"`
	UnlimitedQuota bool     `json:"unlimited_quota" gorm:"default:false"`
	ModelLimitsEnabled bool `json:"model_limits_enabled" gorm:"default:false"`
	ModelLimits   string    `json:"model_limits" gorm:"type:text"`
	UsedQuota     int       `json:"used_quota" gorm:"default:0"`
	Group         string    `json:"group" gorm:"default:default;size:64"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetTokenByKey returns a token by key
func GetTokenByKey(key string) (*Token, error) {
	var token Token
	err := DB.Where("`key` = ?", key).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokenById returns a token by ID
func GetTokenById(id int) (*Token, error) {
	var token Token
	err := DB.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// CreateToken creates a new token
func CreateToken(token *Token) error {
	return DB.Create(token).Error
}

// UpdateToken updates a token
func UpdateToken(token *Token) error {
	return DB.Save(token).Error
}

// DeleteToken deletes a token by ID
func DeleteToken(id int) error {
	return DB.Delete(&Token{}, id).Error
}

// GetTokensByUserId returns tokens by user ID
func GetTokensByUserId(userId int, page, pageSize int) ([]Token, int64, error) {
	var tokens []Token
	var total int64

	err := DB.Model(&Token{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Where("user_id = ?", userId).Offset((page - 1) * pageSize).Limit(pageSize).Find(&tokens).Error
	if err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

// ValidateToken validates a token and returns the associated user
func ValidateToken(key string) (*User, *Token, error) {
	token, err := GetTokenByKey(key)
	if err != nil {
		return nil, nil, err
	}

	// Check token status
	if token.Status != 1 {
		return nil, nil, ErrTokenDisabled
	}

	// Check expiration
	if token.ExpiredTime > 0 && token.ExpiredTime < time.Now().Unix() {
		return nil, nil, ErrTokenExpired
	}

	// Check quota
	if !token.UnlimitedQuota && token.RemainQuota <= 0 {
		return nil, nil, ErrTokenQuotaExceeded
	}

	// Get user
	user, err := GetUserById(token.UserId)
	if err != nil {
		return nil, nil, err
	}

	// Check user status
	if user.Status != 1 {
		return nil, nil, ErrUserDisabled
	}

	return user, token, nil
}
