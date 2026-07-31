package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/stuntdouble/cli/pkg/dlp"
)

const (
	DummyAnthropicKey = "sd-dummy-anthropic-key-000"
	DummyOpenAIKey    = "sd-dummy-openai-key-000"
)

// ZeroTrustProxy intercepts outgoing API calls from sandboxed agents, replacing
// dummy credentials with real host keys on egress so agents never hold real secrets.
type ZeroTrustProxy struct {
	server        *http.Server
	listener      net.Listener
	listenAddr    string
	realAnthropic string
	realOpenAI     string
	dlpScanner    *dlp.Scanner
	mu            sync.RWMutex
}

// NewZeroTrustProxy initializes a local API key substitution proxy
func NewZeroTrustProxy() *ZeroTrustProxy {
	return &ZeroTrustProxy{
		realAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
		realOpenAI:    os.Getenv("OPENAI_API_KEY"),
		dlpScanner:    dlp.NewScanner(),
	}
}

// Start binds to an available loopback port and starts serving requests
func (p *ZeroTrustProxy) Start(ctx context.Context) (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to bind proxy listener: %w", err)
	}

	p.listener = listener
	p.listenAddr = listener.Addr().String()

	p.server = &http.Server{
		Handler: http.HandlerFunc(p.handleProxy),
	}

	go func() {
		if err := p.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("[ZeroTrustProxy] Server error: %v", err)
		}
	}()

	return p.listenAddr, nil
}

// Stop gracefully shuts down the proxy server
func (p *ZeroTrustProxy) Stop() error {
	if p.server != nil {
		return p.server.Shutdown(context.Background())
	}
	return nil
}

// GetListenAddr returns the assigned IP:Port of the proxy
func (p *ZeroTrustProxy) GetListenAddr() string {
	return p.listenAddr
}

// handleProxy performs dummy key substitution and forwards API requests
func (p *ZeroTrustProxy) handleProxy(w http.ResponseWriter, r *http.Request) {
	// 1. Prepare forwarded request
	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	outReq.Header = r.Header.Clone()

	// 2. Perform zero-trust secret substitution
	p.substituteSecret(outReq.Header, "x-api-key", DummyAnthropicKey, p.realAnthropic)
	p.substituteSecret(outReq.Header, "Authorization", DummyAnthropicKey, p.realAnthropic)
	p.substituteSecret(outReq.Header, "Authorization", DummyOpenAIKey, p.realOpenAI)

	// 3. Execute HTTP request
	client := &http.Client{}
	resp, err := client.Do(outReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Proxy egress error: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 4. Copy response back to caller
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *ZeroTrustProxy) substituteSecret(headers http.Header, headerName, dummyValue, realValue string) {
	if realValue == "" {
		return
	}
	val := headers.Get(headerName)
	if strings.Contains(val, dummyValue) {
		headers.Set(headerName, strings.ReplaceAll(val, dummyValue, realValue))
	}
}
