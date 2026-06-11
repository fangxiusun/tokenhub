package i18n

import (
	"embed"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
)

// Init initializes the i18n system
func Init() error {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load translation files
	_, err := bundle.LoadMessageFileFS(localeFS, "locales/zh-CN.yaml")
	if err != nil {
		log.Printf("Warning: Failed to load zh-CN translations: %v", err)
	}

	_, err = bundle.LoadMessageFileFS(localeFS, "locales/en.yaml")
	if err != nil {
		log.Printf("Warning: Failed to load en translations: %v", err)
	}

	// Initialize localizers
	localizers = make(map[string]*i18n.Localizer)
	localizers["zh-CN"] = i18n.NewLocalizer(bundle, "zh-CN")
	localizers["en"] = i18n.NewLocalizer(bundle, "en")

	log.Println("i18n initialized successfully")
	return nil
}

// T translates a message key
func T(lang, key string, args ...map[string]interface{}) string {
	localizer, ok := localizers[lang]
	if !ok {
		localizer = localizers["en"]
	}

	config := &i18n.LocalizeConfig{
		MessageID: key,
	}

	if len(args) > 0 {
		config.TemplateData = args[0]
	}

	msg, err := localizer.Localize(config)
	if err != nil {
		return key
	}
	return msg
}

