package model

import (
	"time"
)

// Channel represents an API channel/provider
type Channel struct {
	Id            int       `json:"id" gorm:"primaryKey"`
	Type          int       `json:"type" gorm:"default:0"`
	Key           string    `json:"key" gorm:"type:text"`
	BaseURL       string    `json:"base_url" gorm:"size:256"`
	Other         string    `json:"other" gorm:"type:text"`
	Status        int       `json:"status" gorm:"default:1"`
	Name          string    `json:"name" gorm:"size:64"`
	Weight        int       `json:"weight" gorm:"default:1"`
	Balance       float64   `json:"balance" gorm:"default:0"`
	BalanceUpdatedTime int64 `json:"balance_updated_time" gorm:"default:0"`
	Models        string    `json:"models" gorm:"type:text"`
	ModelMapping  string    `json:"model_mapping" gorm:"type:text"`
	Groups        string    `json:"groups" gorm:"default:'default';size:256"`
	Tags          string    `json:"tags" gorm:"type:text"`
	AutoBan       bool      `json:"auto_ban" gorm:"default:false"`
	ModelRatio    float64   `json:"model_ratio" gorm:"default:1"`
	ModelPrice    float64   `json:"model_price" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetChannelById returns a channel by ID
func GetChannelById(id int) (*Channel, error) {
	var channel Channel
	err := DB.First(&channel, id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// CreateChannel creates a new channel
func CreateChannel(channel *Channel) error {
	return DB.Create(channel).Error
}

// UpdateChannel updates a channel
func UpdateChannel(channel *Channel) error {
	return DB.Save(channel).Error
}

// DeleteChannel deletes a channel by ID
func DeleteChannel(id int) error {
	return DB.Delete(&Channel{}, id).Error
}

// GetAllChannels returns all channels
func GetAllChannels() ([]Channel, error) {
	var channels []Channel
	err := DB.Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// GetChannelsByStatus returns channels by status
func GetChannelsByStatus(status int) ([]Channel, error) {
	var channels []Channel
	err := DB.Where("status = ?", status).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// GetChannelsByType returns channels by type
func GetChannelsByType(channelType int) ([]Channel, error) {
	var channels []Channel
	err := DB.Where("type = ?", channelType).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}
