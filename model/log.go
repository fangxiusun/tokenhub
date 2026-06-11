package model

import (
	"time"
)

// Log represents an API request log
type Log struct {
	Id               int       `json:"id" gorm:"primaryKey"`
	UserId           int       `json:"user_id" gorm:"index"`
	CreatedAt        time.Time `json:"created_at"`
	Type             int       `json:"type" gorm:"default:0"`
	Content          string    `json:"content" gorm:"type:text"`
	ModelName        string    `json:"model_name" gorm:"size:128"`
	TokenName        string    `json:"token_name" gorm:"size:64"`
	Quota            int       `json:"quota" gorm:"default:0"`
	TokenId          int       `json:"token_id" gorm:"index"`
	Username         string    `json:"username" gorm:"size:64"`
	ChannelId        int       `json:"channel_id" gorm:"index"`
	Channel          string    `json:"channel" gorm:"size:64"`
	PromptTokens     int       `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int       `json:"completion_tokens" gorm:"default:0"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	IsStream         bool      `json:"is_stream" gorm:"default:false"`
	AccessToken      string    `json:"access_token" gorm:"size:64"`
	IPAddress        string    `json:"ip_address" gorm:"size:64"`
	GeoIP            string    `json:"geo_ip" gorm:"size:64"`
}

// CreateLog creates a new log entry
func CreateLog(log *Log) error {
	return DB.Create(log).Error
}

// GetLogsByUserId returns logs by user ID
func GetLogsByUserId(userId int, page, pageSize int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	err := DB.Model(&Log{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Where("user_id = ?", userId).Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAllLogs returns all logs
func GetAllLogs(page, pageSize int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	err := DB.Model(&Log{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DeleteOldLogs deletes logs older than the specified days
func DeleteOldLogs(days int) error {
	return DB.Where("created_at < ?", time.Now().AddDate(0, 0, -days)).Delete(&Log{}).Error
}
