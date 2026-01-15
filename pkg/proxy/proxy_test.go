package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ismailtsdln/interceptify/pkg/ca"
)

func TestProxyIntegration(t *testing.T) {
	// 1. Setup CA
	caCertPath := "test_proxy_ca.crt"
	caKeyPath := "test_proxy_ca.key"
	defer os.Remove(caCertPath)
	defer os.Remove(caKeyPath)

	caInstance, err := ca.NewCA(caCertPath, caKeyPath)
	if err != nil {
		t.Fatalf("failed to setup CA: %v", err)
	}

	// 2. Setup Backend Server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from Backend")
	}))
	defer backend.Close()

	// 3. Setup Proxy
	proxyAddr := "127.0.0.1:9091"
	p := NewProxy(proxyAddr, caInstance)
	go func() {
		if err := p.Start(); err != nil {
			fmt.Printf("Proxy start error: %v\n", err)
		}
	}()

	// Wait for proxy to start
	time.Sleep(100 * time.Millisecond)

	// 4. Test HTTP Request through Proxy
	proxyURL, _ := url.Parse("http://" + proxyAddr)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	resp, err := client.Get(backend.URL)
	if err != nil {
		t.Fatalf("failed to send request through proxy: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello from Backend" {
		t.Errorf("expected 'Hello from Backend', got '%s'", string(body))
	}

	// 5. Test HTTPS (Simulated)
	// For full HTTPS test, we'd need to trust the CA in the client
	// but we can at least check if the proxy handles CONNECT
}
