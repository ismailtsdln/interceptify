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
	"sync"

	"github.com/ismailtsdln/interceptify/pkg/ca"
	"github.com/ismailtsdln/interceptify/pkg/plugins"
	"golang.org/x/net/http2"
)

// Proxy represents the MITM proxy server
type Proxy struct {
	Addr      string
	CA        *ca.CA
	Plugins   *plugins.Manager
	EventChan chan string
	mu        sync.Mutex
	clients   map[chan string]bool
}

// NewProxy creates a new Proxy instance
func NewProxy(addr string, caInstance *ca.CA) *Proxy {
	return &Proxy{
		Addr:      addr,
		CA:        caInstance,
		Plugins:   plugins.NewManager(),
		EventChan: make(chan string, 100),
		clients:   make(map[chan string]bool),
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

	go p.broadcastEvents()

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
	} else if strings.HasPrefix(req.Host, "interceptify.local") || req.Host == "interceptify" || req.Host == "localhost:8080" {
		p.handleDashboard(conn, req)
	} else {
		p.handleHTTP(conn, req)
	}
}

func (p *Proxy) broadcastEvents() {
	for event := range p.EventChan {
		p.mu.Lock()
		for clientChan := range p.clients {
			select {
			case clientChan <- event:
			default:
				// Buffer full, skip or handle
			}
		}
		p.mu.Unlock()
	}
}

func (p *Proxy) logEvent(event string) {
	select {
	case p.EventChan <- event:
	default:
	}
}

func (p *Proxy) handleDashboard(conn net.Conn, req *http.Request) {
	if req.URL.Path == "/events" {
		p.handleSSE(conn)
		return
	}

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Interceptify Dashboard</title>
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600&display=swap" rel="stylesheet">
		<style>
			:root {
				--bg: #0a0a0c;
				--card-bg: rgba(255, 255, 255, 0.03);
				--primary: #00ffcc;
				--secondary: #7000ff;
				--text: #e0e0e0;
				--border: rgba(255, 255, 255, 0.1);
			}
			body {
				font-family: 'Outfit', sans-serif;
				background-color: var(--bg);
				color: var(--text);
				margin: 0;
				padding: 2rem;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
				background-image: 
					radial-gradient(circle at 10% 20%, rgba(112, 0, 255, 0.05) 0%, transparent 40%),
					radial-gradient(circle at 90% 80%, rgba(0, 255, 204, 0.05) 0%, transparent 40%);
			}
			.container {
				max-width: 1000px;
				width: 100%;
			}
			header {
				display: flex;
				justify-content: space-between;
				align-items: center;
				margin-bottom: 3rem;
				width: 100%;
			}
			h1 {
				font-size: 2.5rem;
				font-weight: 600;
				margin: 0;
				background: linear-gradient(to right, var(--primary), var(--secondary));
				-webkit-background-clip: text;
				-webkit-text-fill-color: transparent;
			}
			.status-pill {
				background: rgba(0, 255, 204, 0.1);
				color: var(--primary);
				padding: 0.5rem 1.2rem;
				border-radius: 2rem;
				font-size: 0.9rem;
				font-weight: 600;
				border: 1px solid rgba(0, 255, 204, 0.2);
				display: flex;
				align-items: center;
				gap: 0.5rem;
			}
			.status-pill::before {
				content: '';
				width: 8px;
				height: 8px;
				background: var(--primary);
				border-radius: 50%;
				box-shadow: 0 0 10px var(--primary);
			}
			.grid {
				display: grid;
				grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
				gap: 1.5rem;
				margin-bottom: 2rem;
			}
			.card {
				background: var(--card-bg);
				backdrop-filter: blur(10px);
				border: 1px solid var(--border);
				border-radius: 1.2rem;
				padding: 1.5rem;
				transition: transform 0.3s ease;
			}
			.card:hover {
				transform: translateY(-5px);
				border-color: rgba(255, 255, 255, 0.2);
			}
			.card h3 {
				margin-top: 0;
				font-size: 1.1rem;
				color: rgba(255, 255, 255, 0.6);
			}
			.card .value {
				font-size: 2.5rem;
				font-weight: 600;
				margin: 0.5rem 0;
			}
			.traffic-log {
				background: var(--card-bg);
				backdrop-filter: blur(10px);
				border: 1px solid var(--border);
				border-radius: 1.2rem;
				width: 100%;
				height: 400px;
				overflow-y: auto;
				padding: 1rem;
				display: flex;
				flex-direction: column;
				gap: 0.5rem;
			}
			.log-entry {
				padding: 0.8rem 1rem;
				background: rgba(255, 255, 255, 0.02);
				border-radius: 0.8rem;
				font-family: 'Courier New', Courier, monospace;
				font-size: 0.9rem;
				border-left: 3px solid var(--secondary);
				animation: slideIn 0.3s ease-out;
			}
			@keyframes slideIn {
				from { opacity: 0; transform: translateX(-10px); }
				to { opacity: 1; transform: translateX(0); }
			}
			.method {
				font-weight: bold;
				color: var(--primary);
				margin-right: 10px;
			}
			::-webkit-scrollbar {
				width: 8px;
			}
			::-webkit-scrollbar-track {
				background: transparent;
			}
			::-webkit-scrollbar-thumb {
				background: rgba(255, 255, 255, 0.1);
				border-radius: 4px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<header>
				<h1>Interceptify 📈</h1>
				<div class="status-pill">PROXY ACTIVE</div>
			</header>

			<div class="grid">
				<div class="card">
					<h3>Current Connections</h3>
					<div class="value" id="conn-count">1</div>
				</div>
				<div class="card">
					<h3>Requests Intercepted</h3>
					<div class="value" id="req-count">0</div>
				</div>
			</div>

			<div class="traffic-log" id="log">
				<!-- Logs will appear here -->
			</div>
		</div>

		<script>
			const logEl = document.getElementById('log');
			const reqCountEl = document.getElementById('req-count');
			let reqCount = 0;

			const eventSource = new EventSource('/events');
			eventSource.onmessage = (event) => {
				const entry = document.createElement('div');
				entry.className = 'log-entry';
				entry.textContent = event.data;
				logEl.prepend(entry);
				
				reqCount++;
				reqCountEl.textContent = reqCount;

				if (logEl.children.length > 50) {
					logEl.removeChild(logEl.lastChild);
				}
			};
		</script>
	</body>
	</html>
	`
	resp := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(html)),
	}
	resp.Header.Set("Content-Type", "text/html")
	resp.Write(conn)
}

func (p *Proxy) handleSSE(conn net.Conn) {
	// Upgrade connection to SSE
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/event-stream\r\n"))
	conn.Write([]byte("Cache-Control: no-cache\r\n"))
	conn.Write([]byte("Connection: keep-alive\r\n\r\n"))

	clientChan := make(chan string, 10)
	p.mu.Lock()
	p.clients[clientChan] = true
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.clients, clientChan)
		p.mu.Unlock()
		close(clientChan)
	}()

	for event := range clientChan {
		_, err := fmt.Fprintf(conn, "data: %s\n\n", event)
		if err != nil {
			return
		}
	}
}

func (p *Proxy) handleHTTP(conn net.Conn, req *http.Request) {
	log.Printf("HTTP Request: %s %s", req.Method, req.URL.String())
	p.logEvent(fmt.Sprintf("HTTP: %s %s", req.Method, req.URL.String()))

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
	p.logEvent(fmt.Sprintf("CONNECT: %s", req.Host))

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
		NextProtos:   []string{"h2", "http/1.1"},
	}

	tlsConn := tls.Server(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed for %s: %v", host, err)
		return
	}
	defer tlsConn.Close()

	// Check negotiated protocol
	if tlsConn.ConnectionState().NegotiatedProtocol == "h2" {
		log.Printf("HTTP/2 Negotiated for %s", host)
		p.handleHTTPS2(tlsConn, req.Host)
		return
	}

	// Now we have a decrypted stream (tlsConn)
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

func (p *Proxy) handleHTTPS2(conn net.Conn, host string) {
	s2 := &http2.Server{}
	s2.ServeConn(conn, &http2.ServeConnOpts{
		Context: nil,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Construct the full URL for H2
			r.URL.Scheme = "https"
			r.URL.Host = host
			p.handleInterceptedRequestH2(w, r)
		}),
	})
}

func (p *Proxy) handleInterceptedRequestH2(w http.ResponseWriter, req *http.Request) {
	// Run Request Hooks
	req, shortCircuitResp := p.Plugins.RunRequestHooks(req)
	if shortCircuitResp != nil {
		copyResponse(w, shortCircuitResp)
		return
	}

	log.Printf("Intercepted HTTPS/2 Request: %s %s", req.Method, req.URL.String())
	p.logEvent(fmt.Sprintf("HTTPS/2: %s %s", req.Method, req.URL.String()))

	// Implement forwarding logic for HTTPS/2
	req.RequestURI = ""
	transport := &http2.Transport{}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to forward intercepted H2 request: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Run Response Hooks
	resp = p.Plugins.RunResponseHooks(req, resp)

	copyResponse(w, resp)
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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
