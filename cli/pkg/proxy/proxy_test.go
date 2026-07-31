package proxy

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestZeroTrustProxySubstitution(t *testing.T) {
	// Setup target mock API server
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer real-openai-secret-999" {
			t.Errorf("Expected real key in Authorization header, got: %s", authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"authenticated"}`))
	}))
	defer targetServer.Close()

	os.Setenv("OPENAI_API_KEY", "real-openai-secret-999")
	defer os.Unsetenv("OPENAI_API_KEY")

	proxy := NewZeroTrustProxy()
	ctx := context.Background()
	_, err := proxy.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start proxy: %v", err)
	}
	defer proxy.Stop()

	// Client sends request to proxy with dummy token targeting mock API
	req, err := http.NewRequest("GET", targetServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+DummyOpenAIKey)

	// Send request through proxy handler directly
	rec := httptest.NewRecorder()
	proxy.handleProxy(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected HTTP 200 OK from proxy, got: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"authenticated"}` {
		t.Fatalf("Unexpected response body: %s", string(body))
	}
}
