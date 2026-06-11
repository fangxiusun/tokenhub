package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PasskeyCredential stores WebAuthn/Passkey credentials
type PasskeyCredential struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	UserId       int       `json:"user_id" gorm:"uniqueIndex"`
	CredentialID string    `json:"credential_id" gorm:"type:varchar(512);uniqueIndex"`
	PublicKey    string    `json:"public_key" gorm:"type:text"`
	SignCount    uint32    `json:"sign_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetPasskeyByUserId returns passkey credential by user ID
func GetPasskeyByUserId(userId int) (*PasskeyCredential, error) {
	var cred PasskeyCredential
	err := DB.Where("user_id = ?", userId).First(&cred).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cred, nil
}

// GetPasskeyByCredentialId returns passkey credential by credential ID
func GetPasskeyByCredentialId(credentialId string) (*PasskeyCredential, error) {
	var cred PasskeyCredential
	err := DB.Where("credential_id = ?", credentialId).First(&cred).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cred, nil
}

// CreatePasskeyCredential creates a new passkey credential
func CreatePasskeyCredential(cred *PasskeyCredential) error {
	return DB.Create(cred).Error
}

// UpdatePasskeyCredential updates a passkey credential
func UpdatePasskeyCredential(cred *PasskeyCredential) error {
	return DB.Save(cred).Error
}

// DeletePasskeyByUserId deletes passkey credential for a user
func DeletePasskeyByUserId(userId int) error {
	return DB.Where("user_id = ?", userId).Delete(&PasskeyCredential{}).Error
}

// ToCredentialBytes returns credential ID and public key as bytes
func (p *PasskeyCredential) ToCredentialBytes() (credentialID, publicKey []byte, err error) {
	credentialID, err = base64.StdEncoding.DecodeString(p.CredentialID)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return credentialID, publicKey, nil
}

// TransportList returns transport list
func (p *PasskeyCredential) TransportList() []string {
	return []string{"internal"}
}

// GetTransportsJSON returns transports as JSON
func (p *PasskeyCredential) GetTransportsJSON() string {
	transports := []string{"internal"}
	data, _ := json.Marshal(transports)
	return string(data)
}
