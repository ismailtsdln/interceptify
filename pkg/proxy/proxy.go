package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/ismailtsdln/interceptify/pkg/ca"
	"github.com/ismailtsdln/interceptify/pkg/plugins"
)

// Proxy represents the MITM proxy server
type Proxy struct {
	Addr    string
	CA      *ca.CA
	Plugins *plugins.Manager
}

// NewProxy creates a new Proxy instance
func NewProxy(addr string, caInstance *ca.CA) *Proxy {
	return &Proxy{
		Addr:    addr,
		CA:      caInstance,
		Plugins: plugins.NewManager(),
	}
}

// Start runs the proxy server
func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", p.Addr, err)
	}
	defer listener.Close()

	log.Printf("Interceptify proxy listening on %s", p.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go p.handleConnection(conn)
	}
}

func (p *Proxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// peak at the first few bytes to decide if it's HTTP or something else
	// or just use http.ReadRequest if we expect HTTP proxy behavior
	req, err := http.ReadRequest(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("failed to read request: %v", err)
		}
		return
	}

	if req.Method == http.MethodConnect {
		p.handleHTTPS(conn, req)
	} else {
		p.handleHTTP(conn, req)
	}
}

func (p *Proxy) handleHTTP(conn net.Conn, req *http.Request) {
	log.Printf("HTTP Request: %s %s", req.Method, req.URL.String())

	// Simple transparent proxy logic or explicit proxy logic
	// For now, just forward and log

	client := &http.Client{}

	// Clean up request for forwarding
	req.RequestURI = ""

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to forward request: %v", err)
		return
	}
	defer resp.Body.Close()

	resp.Write(conn)
}

func (p *Proxy) handleHTTPS(conn net.Conn, req *http.Request) {
	log.Printf("HTTPS Tunnel Request: %s", req.Host)

	// Acknowledge the CONNECT request
	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Strip port from host
	host := strings.Split(req.Host, ":")[0]

	// Sign a certificate for this host
	cert, key, err := p.CA.SignCertificate(host)
	if err != nil {
		log.Printf("failed to sign certificate for %s: %v", host, err)
		return
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	tlsConn := tls.Server(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed for %s: %v", host, err)
		return
	}
	defer tlsConn.Close()

	// Now we have a decrypted stream (tlsConn)
	// We can read HTTP requests from it
	tlsReader := bufio.NewReader(tlsConn)
	for {
		interceptedReq, err := http.ReadRequest(tlsReader)
		if err != nil {
			if err != io.EOF {
				log.Printf("failed to read intercepted request: %v", err)
			}
			break
		}

		p.handleInterceptedRequest(tlsConn, interceptedReq)
	}
}

func (p *Proxy) handleInterceptedRequest(conn net.Conn, req *http.Request) {
	// Run Request Hooks
	req, shortCircuitResp := p.Plugins.RunRequestHooks(req)
	if shortCircuitResp != nil {
		shortCircuitResp.Write(conn)
		return
	}

	log.Printf("Intercepted HTTPS Request: %s %s", req.Method, req.URL.String())

	// Implement forwarding logic for HTTPS
	destURL := "https://" + req.Host + req.URL.String()
	proxyReq, err := http.NewRequest(req.Method, destURL, req.Body)
	if err != nil {
		log.Printf("failed to create proxy request: %v", err)
		return
	}

	// Copy headers
	for k, v := range req.Header {
		proxyReq.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("failed to forward intercepted request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Run Response Hooks
	resp = p.Plugins.RunResponseHooks(req, resp)

	resp.Write(conn)
}
