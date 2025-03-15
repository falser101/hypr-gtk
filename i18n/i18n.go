package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	currentLang  = "en"
	translations = make(map[string]map[string]string)
)

// Initialize loads all translation files from the given directory
func Initialize(langDir string) error {
	entries, err := os.ReadDir(langDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			langCode := filepath.Base(entry.Name()[:len(entry.Name())-5]) // remove .json
			data, err := os.ReadFile(filepath.Join(langDir, entry.Name()))
			if err != nil {
				return err
			}

			var langMap map[string]string
			if err := json.Unmarshal(data, &langMap); err != nil {
				return err
			}
			translations[langCode] = langMap
		}
	}
	return nil
}

// SetLanguage changes the current language
func SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		currentLang = lang
	}
}

// GetLanguage returns the current language code
func GetLanguage() string {
	return currentLang
}

// Tr returns the translated string for the given key
func Tr(key string) string {
	if trans, ok := translations[currentLang]; ok {
		if str, ok := trans[key]; ok {
			return str
		}
	}
	return key
}

// GetAvailableLanguages returns a list of available language codes
func GetAvailableLanguages() []string {
	langs := make([]string, 0, len(translations))
	for lang := range translations {
		langs = append(langs, lang)
	}
	return langs
}
