package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	Id            int       `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"uniqueIndex;size:64"`
	Password      string    `json:"-" gorm:"size:128"`
	DisplayName   string    `json:"display_name" gorm:"size:64"`
	Role          int       `json:"role" gorm:"default:1"`
	Status        int       `json:"status" gorm:"default:1"`
	Email         string    `json:"email" gorm:"size:128"`
	Quota         int       `json:"quota" gorm:"default:0"`
	UsedQuota     int       `json:"used_quota" gorm:"default:0"`
	RequestCount  int       `json:"request_count" gorm:"default:0"`
	Group         string    `json:"group" gorm:"default:default;size:64"`
	AffCode       string    `json:"aff_code" gorm:"uniqueIndex;size:16"`
	AffiliatedQuota int     `json:"affiliated_quota" gorm:"default:0"`
	InviterId     int       `json:"inviter_id" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetUserById returns a user by ID
func GetUserById(id int) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername returns a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// UpdateUser updates a user
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser deletes a user by ID
func DeleteUser(id int) error {
	return DB.Delete(&User{}, id).Error
}

// GetAllUsers returns all users
func GetAllUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	err := DB.Model(&User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SearchUsers searches users by keyword
func SearchUsers(keyword string, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	query := DB.Model(&User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
