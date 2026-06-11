package model

// Option represents a system option/configuration
type Option struct {
	Key   string `json:"key" gorm:"primaryKey;size:64"`
	Value string `json:"value" gorm:"type:text"`
}

// GetOption returns an option by key
func GetOption(key string) (string, error) {
	var option Option
	err := DB.Where("`key` = ?", key).First(&option).Error
	if err != nil {
		return "", err
	}
	return option.Value, nil
}

// SetOption sets an option value
func SetOption(key, value string) error {
	option := Option{
		Key:   key,
		Value: value,
	}
	return DB.Save(&option).Error
}

// GetOptionOrDefault returns an option value or a default value
func GetOptionOrDefault(key, defaultValue string) string {
	value, err := GetOption(key)
	if err != nil {
		return defaultValue
	}
	return value
}
