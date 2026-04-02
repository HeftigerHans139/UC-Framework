package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type languageSettingsResponse struct {
	Language           string   `json:"language"`
	DefaultLanguage    string   `json:"default_language"`
	SupportedLanguages []string `json:"supported_languages"`
}

type languageSettingsRequest struct {
	Language string `json:"language"`
}

var GetLanguageSettingsFunc func() (string, []string, error)
var SaveLanguageSettingsFunc func(language string) error

func normalizeLanguageCode(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "de":
		return "de"
	default:
		return "en"
	}
}

func normalizeSupportedLanguages(languages []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(languages))
	for _, language := range languages {
		normalized := normalizeLanguageCode(language)
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	if len(result) == 0 {
		return []string{"en", "de"}
	}
	return result
}

func isSupportedLanguage(language string, supported []string) bool {
	for _, candidate := range supported {
		if candidate == language {
			return true
		}
	}
	return false
}

func writeLanguageSettingsResponse(w http.ResponseWriter, language string, supported []string) {
	response := languageSettingsResponse{
		Language:           language,
		DefaultLanguage:    language,
		SupportedLanguages: supported,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func LanguageSettingsHandler(w http.ResponseWriter, r *http.Request) {
	language := "en"
	supported := []string{"en", "de"}
	if GetLanguageSettingsFunc != nil {
		storedLanguage, storedSupported, err := GetLanguageSettingsFunc()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load language settings: %v", err), http.StatusInternalServerError)
			return
		}
		language = normalizeLanguageCode(storedLanguage)
		supported = normalizeSupportedLanguages(storedSupported)
	}

	switch r.Method {
	case http.MethodGet:
		writeLanguageSettingsResponse(w, language, supported)

	case http.MethodPost:
		var req languageSettingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		nextLanguage := normalizeLanguageCode(req.Language)
		if !isSupportedLanguage(nextLanguage, supported) {
			http.Error(w, "unsupported language", http.StatusBadRequest)
			return
		}
		if SaveLanguageSettingsFunc == nil {
			http.Error(w, "language settings unavailable", http.StatusServiceUnavailable)
			return
		}
		if err := SaveLanguageSettingsFunc(nextLanguage); err != nil {
			http.Error(w, fmt.Sprintf("failed to save language settings: %v", err), http.StatusInternalServerError)
			return
		}
		writeLanguageSettingsResponse(w, nextLanguage, supported)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
