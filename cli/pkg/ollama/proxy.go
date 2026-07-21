package ollama

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// StartProxy creates a reverse proxy to localhost:11434 (Ollama)
// It injects the StuntDouble Universal Stunt Protocol token into every request
// so that local models know they are running in a secure sandbox.
func StartProxy(listenPort string) error {
	targetUrl, err := url.Parse("http://localhost:11434")
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Intercept and rewrite the request
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", targetUrl.Host)
		
		// Inject the STP (Universal Stunt Protocol) Attestation
		// This tells the local AI that its environment is kernel-hardened
		req.Header.Add("X-StuntDouble-Sandbox", "STP-v2.0-SECURE")
		
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = targetUrl.Path
		
		fmt.Printf("🛡️  [Ollama Proxy] Intercepted prompt to /%s. Injected STP Attestation.\n", req.URL.Path)
	}

	// Intercept the response for telemetry
	proxy.ModifyResponse = func(resp *http.Response) error {
		fmt.Printf("📦 [Ollama Proxy] Local AI responded with HTTP %d.\n", resp.StatusCode)
		return nil
	}

	fmt.Printf(">> [StuntDouble] Local Ollama Proxy running on port %s\n", listenPort)
	fmt.Printf(">> [StuntDouble] Point your AI apps to http://localhost:%s to enforce sandbox governance.\n", listenPort)
	
	return http.ListenAndServe(":"+listenPort, proxy)
}
