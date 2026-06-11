package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// TwoFactorAuth stores 2FA configuration for users
type TwoFactorAuth struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	UserId    int       `json:"user_id" gorm:"uniqueIndex"`
	Secret    string    `json:"-" gorm:"size:128"`
	Enabled   bool      `json:"enabled" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetTwoFactorAuthByUserId returns 2FA config by user ID
func GetTwoFactorAuthByUserId(userId int) (*TwoFactorAuth, error) {
	var tfa TwoFactorAuth
	err := DB.Where("user_id = ?", userId).First(&tfa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tfa, nil
}

// CreateOrUpdateTwoFactorAuth creates or updates 2FA config
func CreateOrUpdateTwoFactorAuth(tfa *TwoFactorAuth) error {
	var existing TwoFactorAuth
	err := DB.Where("user_id = ?", tfa.UserId).First(&existing).Error
	if err == nil {
		// Update existing
		existing.Secret = tfa.Secret
		existing.Enabled = tfa.Enabled
		return DB.Save(&existing).Error
	}
	// Create new
	return DB.Create(tfa).Error
}

// IsTwoFAEnabled checks if 2FA is enabled for a user
func IsTwoFAEnabled(userId int) bool {
	var count int64
	DB.Model(&TwoFactorAuth{}).Where("user_id = ? AND enabled = ?", userId, true).Count(&count)
	return count > 0
}

// DeleteTwoFactorAuth deletes 2FA config for a user
func DeleteTwoFactorAuth(userId int) error {
	return DB.Where("user_id = ?", userId).Delete(&TwoFactorAuth{}).Error
}

