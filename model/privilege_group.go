package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PrivilegeGroup represents a privilege group (e.g., svip, vip, default)
type PrivilegeGroup struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:64"`
	DisplayName string    `json:"display_name" gorm:"size:128"`
	Description string    `json:"description" gorm:"size:256"`
	Quota       int       `json:"quota" gorm:"default:0"`
	RateLimit   int       `json:"rate_limit" gorm:"default:0"`
	Models      string    `json:"models" gorm:"type:text"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// GetPrivilegeGroupById returns a privilege group by ID
func GetPrivilegeGroupById(id int) (*PrivilegeGroup, error) {
	var group PrivilegeGroup
	err := DB.First(&group, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

// GetPrivilegeGroupByName returns a privilege group by name
func GetPrivilegeGroupByName(name string) (*PrivilegeGroup, error) {
	var group PrivilegeGroup
	err := DB.Where("name = ?", name).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

// GetAllPrivilegeGroups returns all privilege groups
func GetAllPrivilegeGroups() ([]PrivilegeGroup, error) {
	var groups []PrivilegeGroup
	err := DB.Order("id asc").Find(&groups).Error
	return groups, err
}

// CreatePrivilegeGroup creates a new privilege group
func CreatePrivilegeGroup(group *PrivilegeGroup) error {
	return DB.Create(group).Error
}

// UpdatePrivilegeGroup updates a privilege group
func UpdatePrivilegeGroup(group *PrivilegeGroup) error {
	return DB.Save(group).Error
}

// DeletePrivilegeGroup soft-deletes a privilege group
func DeletePrivilegeGroup(id int) error {
	return DB.Delete(&PrivilegeGroup{}, id).Error
}

// GetDefaultPrivilegeGroup returns the default privilege group
func GetDefaultPrivilegeGroup() (*PrivilegeGroup, error) {
	return GetPrivilegeGroupByName("default")
}

// InitializeDefaultPrivilegeGroup creates the default group if it doesn't exist
func InitializeDefaultPrivilegeGroup() error {
	var count int64
	DB.Model(&PrivilegeGroup{}).Count(&count)
	if count > 0 {
		return nil
	}

	defaultGroup := PrivilegeGroup{
		Name:        "default",
		DisplayName: "Default",
		Description: "Default privilege group",
		Quota:       100000,
		RateLimit:   60,
		Status:      1,
	}
	return DB.Create(&defaultGroup).Error
}

