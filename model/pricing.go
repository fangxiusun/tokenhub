package model

// Pricing represents model pricing information
type Pricing struct {
	ModelName       string  `json:"model_name" gorm:"primaryKey;size:128"`
	QuotaType       int     `json:"quota_type" gorm:"default:0"`
	ModelRatio      float64 `json:"model_ratio" gorm:"default:1"`
	ModelPrice      float64 `json:"model_price" gorm:"default:0"`
	CompletionRatio float64 `json:"completion_ratio" gorm:"default:1"`
	CacheRatio      float64 `json:"cache_ratio" gorm:"default:0.5"`
	AudioRatio      float64 `json:"audio_ratio" gorm:"default:1"`
}

// GetPricingByModel returns pricing by model name
func GetPricingByModel(modelName string) (*Pricing, error) {
	var pricing Pricing
	err := DB.Where("model_name = ?", modelName).First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

// GetAllPricing returns all pricing information
func GetAllPricing() ([]Pricing, error) {
	var pricing []Pricing
	err := DB.Find(&pricing).Error
	if err != nil {
		return nil, err
	}
	return pricing, nil
}

// CreatePricing creates new pricing information
func CreatePricing(pricing *Pricing) error {
	return DB.Create(pricing).Error
}

// UpdatePricing updates pricing information
func UpdatePricing(pricing *Pricing) error {
	return DB.Save(pricing).Error
}

// DeletePricing deletes pricing by model name
func DeletePricing(modelName string) error {
	return DB.Where("model_name = ?", modelName).Delete(&Pricing{}).Error
}

