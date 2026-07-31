package mock

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSyntheticMockGeneratorAndHandler(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stuntdouble-mock-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	gen := NewGenerator(tmpDir)
	mockFile, err := gen.GenerateConfigFile()
	if err != nil {
		t.Fatalf("Failed to generate mock file: %v", err)
	}

	if _, err := os.Stat(mockFile); os.IsNotExist(err) {
		t.Fatalf("Mock file was not created: %s", mockFile)
	}

	handler := NewMockHandler(nil)
	req := httptest.NewRequest("GET", "/v1/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("X-StuntDouble-Synthetic") != "fallback" {
		t.Errorf("Expected fallback header, got: %s", rec.Header().Get("X-StuntDouble-Synthetic"))
	}
}
