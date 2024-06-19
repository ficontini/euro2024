package openliga

import (
	"embed"
	"encoding/json"
)

const translation_file = "translations.json"

type Translator interface {
	Get(string) string
}

//go:embed translations.json
var translationsFile embed.FS

type translator struct {
	translations map[string]string
}

func NewTranslator() (Translator, error) {
	translations, err := loadTranslations()
	if err != nil {
		return nil, err
	}
	return &translator{
		translations: translations,
	}, nil
}

func (t *translator) Get(key string) string {
	if translated, exists := t.translations[key]; exists {
		return translated
	}
	return ""
}

func loadTranslations() (map[string]string, error) {
	bytes, err := translationsFile.ReadFile(translation_file)
	if err != nil {
		return nil, err
	}

	var translations map[string]string
	if err := json.Unmarshal(bytes, &translations); err != nil {
		return nil, err
	}
	return translations, nil

}
