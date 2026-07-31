package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type MockRule struct {
	URLPattern   string            `json:"url_pattern"`
	Method       string            `json:"method"`
	ResponseCode int               `json:"response_code"`
	ResponseBody map[string]any    `json:"response_body"`
	Headers      map[string]string `json:"headers"`
}

type Generator struct {
	workspace string
}

func NewGenerator(workspace string) *Generator {
	return &Generator{workspace: workspace}
}

// GenerateConfigFile creates a default synthetic mock template at .stuntdouble/mocks.json
func (g *Generator) GenerateConfigFile() (string, error) {
	stuntDir := filepath.Join(g.workspace, ".stuntdouble")
	if err := os.MkdirAll(stuntDir, 0755); err != nil {
		return "", err
	}

	mockPath := filepath.Join(stuntDir, "mocks.json")
	defaultMocks := []MockRule{
		{
			URLPattern:   "/v1/charges",
			Method:       "POST",
			ResponseCode: 200,
			ResponseBody: map[string]any{
				"id":       "ch_synthetic_1337",
				"status":   "succeeded",
				"amount":   2000,
				"currency": "usd",
			},
			Headers: map[string]string{"Content-Type": "application/json", "X-StuntDouble-Synthetic": "true"},
		},
		{
			URLPattern:   "/oauth/token",
			Method:       "POST",
			ResponseCode: 200,
			ResponseBody: map[string]any{
				"access_token": "sd_synthetic_access_token_99",
				"token_type":   "Bearer",
				"expires_in":   3600,
			},
			Headers: map[string]string{"Content-Type": "application/json", "X-StuntDouble-Synthetic": "true"},
		},
	}

	data, err := json.MarshalIndent(defaultMocks, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(mockPath, data, 0644); err != nil {
		return "", err
	}

	return mockPath, nil
}

// MockHandler matches requests against synthetic rules
type MockHandler struct {
	rules []MockRule
}

func NewMockHandler(rules []MockRule) *MockHandler {
	return &MockHandler{rules: rules}
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, rule := range m.rules {
		if strings.Contains(r.URL.Path, rule.URLPattern) && (rule.Method == "" || strings.EqualFold(r.Method, rule.Method)) {
			for k, v := range rule.Headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(rule.ResponseCode)
			json.NewEncoder(w).Encode(rule.ResponseBody)
			return
		}
	}

	// Dynamic fallback
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-StuntDouble-Synthetic", "fallback")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": fmt.Sprintf("Synthetic dynamic mock for path: %s", r.URL.Path),
	})
}
