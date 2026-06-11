package model

// Ability represents the mapping between group, model, and channel
type Ability struct {
	Group     string `json:"group" gorm:"primaryKey;size:64"`
	Model     string `json:"model" gorm:"primaryKey;size:128"`
	ChannelId int    `json:"channel_id" gorm:"primaryKey"`
	Enabled   bool   `json:"enabled" gorm:"default:true"`
	Priority  int    `json:"priority" gorm:"default:0"`
	Weight    uint   `json:"weight" gorm:"default:1"`
}

// GetAbilitiesByGroupAndModel returns abilities by group and model
func GetAbilitiesByGroupAndModel(group, model string) ([]Ability, error) {
	var abilities []Ability
	err := DB.Where("`group` = ? AND model = ? AND enabled = ?", group, model, true).
		Order("priority DESC").Find(&abilities).Error
	if err != nil {
		return nil, err
	}
	return abilities, nil
}

// CreateAbility creates a new ability
func CreateAbility(ability *Ability) error {
	return DB.Create(ability).Error
}

// DeleteAbilitiesByChannelId deletes all abilities for a channel
func DeleteAbilitiesByChannelId(channelId int) error {
	return DB.Where("channel_id = ?", channelId).Delete(&Ability{}).Error
}

// GetRandomSatisfiedChannel returns a random channel that satisfies the group and model
func GetRandomSatisfiedChannel(group, model string) (*Channel, error) {
	var abilities []Ability
	err := DB.Where("`group` = ? AND model = ? AND enabled = ?", group, model, true).
		Order("priority DESC, RANDOM()").Find(&abilities).Error
	if err != nil {
		return nil, err
	}

	if len(abilities) == 0 {
		return nil, ErrChannelNotFound
	}

	// Get the first channel (highest priority, random weight)
	return GetChannelById(abilities[0].ChannelId)
}

