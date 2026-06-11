package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/constant"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	Id            int            `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"uniqueIndex;size:64"`
	Password      string         `json:"-" gorm:"size:128"`
	DisplayName   string         `json:"display_name" gorm:"size:64"`
	Role          int            `json:"role" gorm:"default:1"`
	Status        int            `json:"status" gorm:"default:1"`
	Email         string         `json:"email" gorm:"index;size:128"`
	Quota         int            `json:"quota" gorm:"default:0"`
	UsedQuota     int            `json:"used_quota" gorm:"default:0"`
	RequestCount  int            `json:"request_count" gorm:"default:0"`
	Group         string         `json:"group" gorm:"default:default;size:64"`
	AffCode       string         `json:"aff_code" gorm:"uniqueIndex;size:16"`
	AffCount      int            `json:"aff_count" gorm:"default:0"`
	AffQuota      int            `json:"aff_quota" gorm:"default:0"`
	InviterId     int            `json:"inviter_id" gorm:"default:0"`
	AccessToken   *string        `json:"-" gorm:"type:char(32);column:access_token;uniqueIndex"`
	TFAEnabled    bool           `json:"tfa_enabled" gorm:"default:false"`
	PasskeyEnabled bool          `json:"passkey_enabled" gorm:"default:false"`
	Setting       string         `json:"setting" gorm:"type:text"`
	Remark        string         `json:"remark,omitempty" gorm:"size:256"`
	LastLoginAt   int64          `json:"last_login_at" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// UserSetting stores user preferences
type UserSetting struct {
	Language      string `json:"language,omitempty"`
	Theme         string `json:"theme,omitempty"`
	NotifyEmail   bool   `json:"notify_email,omitempty"`
}

// GetUserById returns a user by ID
func GetUserById(id int, selectAll ...bool) (*User, error) {
	var user User
	var err error
	if len(selectAll) > 0 && selectAll[0] {
		err = DB.First(&user, "id = ?", id).Error
	} else {
		err = DB.Omit("password").First(&user, "id = ?", id).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername returns a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail returns a user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// CheckUserExistOrDeleted checks if user exists or was deleted
func CheckUserExistOrDeleted(username, email string) (bool, error) {
	var user User
	var err error
	if email == "" {
		err = DB.Unscoped().First(&user, "username = ?", username).Error
	} else {
		err = DB.Unscoped().First(&user, "username = ? OR email = ?", username, email).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ValidateAndFill validates password and fills user data
func (user *User) ValidateAndFill() error {
	if user.Username == "" || user.Password == "" {
		return ErrUserEmptyCredentials
	}
	password := user.Password
	err := DB.Where("username = ?", user.Username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}
		return err
	}
	if !common.ValidatePasswordAndHash(password, user.Password) {
		return ErrInvalidCredentials
	}
	if user.Status != constant.UserStatusEnabled {
		return ErrInvalidCredentials
	}
	return nil
}

// Insert creates a new user
func (user *User) Insert() error {
	var err error
	if user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	user.AffCode = common.GetRandomString(4)
	return DB.Create(user).Error
}

// Update updates user fields
func (user *User) Update(updatePassword bool) error {
	var err error
	if updatePassword && user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	newUser := *user
	DB.First(&user, user.Id)

	updates := map[string]interface{}{
		"username":       newUser.Username,
		"display_name":   newUser.DisplayName,
		"email":          newUser.Email,
		"group":          newUser.Group,
		"role":           newUser.Role,
		"status":         newUser.Status,
		"quota":          newUser.Quota,
		"used_quota":     newUser.UsedQuota,
		"request_count":  newUser.RequestCount,
		"aff_code":       newUser.AffCode,
		"aff_count":      newUser.AffCount,
		"aff_quota":      newUser.AffQuota,
		"inviter_id":     newUser.InviterId,
		"access_token":   newUser.AccessToken,
		"tfa_enabled":    newUser.TFAEnabled,
		"passkey_enabled": newUser.PasskeyEnabled,
		"setting":        newUser.Setting,
		"remark":         newUser.Remark,
		"last_login_at":  newUser.LastLoginAt,
	}
	if updatePassword {
		updates["password"] = newUser.Password
	}
	return DB.Model(user).Updates(updates).Error
}

// Edit updates specific user fields
func (user *User) Edit(updatePassword bool) error {
	var err error
	if updatePassword && user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	newUser := *user
	updates := map[string]interface{}{
		"username":     newUser.Username,
		"display_name": newUser.DisplayName,
		"group":        newUser.Group,
	}
	if updatePassword {
		updates["password"] = newUser.Password
	}
	DB.First(&user, user.Id)
	return DB.Model(user).Updates(updates).Error
}

// Delete soft-deletes a user
func (user *User) Delete() error {
	if user.Id == 0 {
		return errors.New("id is empty")
	}
	return DB.Delete(user).Error
}

// HardDelete permanently deletes a user
func (user *User) HardDelete() error {
	if user.Id == 0 {
		return errors.New("id is empty")
	}
	return DB.Unscoped().Delete(user).Error
}

// GetAccessToken returns the access token
func (user *User) GetAccessToken() string {
	if user.AccessToken == nil {
		return ""
	}
	return *user.AccessToken
}

// SetAccessToken sets the access token
func (user *User) SetAccessToken(token string) {
	user.AccessToken = &token
}

// GetSetting returns user settings
func (user *User) GetSetting() UserSetting {
	setting := UserSetting{}
	if user.Setting != "" {
		_ = json.Unmarshal([]byte(user.Setting), &setting)
	}
	return setting
}

// SetSetting sets user settings
func (user *User) SetSetting(setting UserSetting) {
	data, _ := json.Marshal(setting)
	user.Setting = string(data)
}

// IsAdmin checks if the user is an admin
func IsAdmin(userId int) bool {
	if userId == 0 {
		return false
	}
	var user User
	err := DB.Where("id = ?", userId).Select("role").Find(&user).Error
	if err != nil {
		return false
	}
	return user.Role >= constant.RoleAdminUser
}

// IsRoot checks if the user is a root user
func IsRoot(userId int) bool {
	if userId == 0 {
		return false
	}
	var user User
	err := DB.Where("id = ?", userId).Select("role").Find(&user).Error
	if err != nil {
		return false
	}
	return user.Role == constant.RoleRootUser
}

// ValidateAccessToken validates an access token
func ValidateAccessToken(token string) (*User, error) {
	if token == "" {
		return nil, nil
	}
	user := &User{}
	err := DB.Where("access_token = ?", token).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserLastLoginAt updates the last login timestamp
func UpdateUserLastLoginAt(id int) {
	DB.Model(&User{}).Where("id = ?", id).Update("last_login_at", time.Now().Unix())
}

// GetAllUsers returns paginated users
func GetAllUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	err := DB.Model(&User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = DB.Omit("password").Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
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
	err = query.Omit("password").Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// IncreaseUserQuota increases user quota
func IncreaseUserQuota(id int, quota int) error {
	return DB.Model(&User{}).Where("id = ?", id).Update("quota", gorm.Expr("quota + ?", quota)).Error
}

// DecreaseUserQuota decreases user quota
func DecreaseUserQuota(id int, quota int) error {
	return DB.Model(&User{}).Where("id = ?", id).Update("quota", gorm.Expr("quota - ?", quota)).Error
}

// UpdateUserUsedQuotaAndRequestCount updates used quota and request count
func UpdateUserUsedQuotaAndRequestCount(id int, quota int) {
	DB.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"used_quota":    gorm.Expr("used_quota + ?", quota),
		"request_count": gorm.Expr("request_count + ?", 1),
	})
}

// RootUserExists checks if a root user exists
func RootUserExists() bool {
	var user User
	err := DB.Where("role = ?", constant.RoleRootUser).First(&user).Error
	return err == nil
}

