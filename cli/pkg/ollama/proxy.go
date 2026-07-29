// Package ollama provides a logging reverse proxy in front of a local Ollama
// server.
//
// The proxy exists for visibility: it records which endpoints an agent calls.
// It does not sandbox the model, block requests, or restrict what the model can
// do, and no header it sets changes model behaviour.
package ollama

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// StartProxy runs a reverse proxy on listenPort that forwards to a local Ollama
// server and logs each request.
func StartProxy(listenPort string) error {
	target, err := url.Parse("http://localhost:11434")
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// NewSingleHostReverseProxy already rewrites scheme, host and path. Wrap its
	// Director rather than replacing it: overwriting req.URL.Path with the
	// target's (empty) path sent every request to "/".
	base := proxy.Director
	proxy.Director = func(req *http.Request) {
		base(req)
		log.Printf("→ %s %s", req.Method, req.URL.Path)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("← %d %s", resp.StatusCode, resp.Request.URL.Path)
		return nil
	}

	addr := "127.0.0.1:" + listenPort
	srv := &http.Server{
		Addr:              addr,
		Handler:           proxy,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Printf("Ollama logging proxy on http://%s → %s\n", addr, target)
	fmt.Println("Requests are logged, not filtered.")
	return srv.ListenAndServe()
}
