<div align="center">
  <img src="assets/banner.png" alt="Interceptify Banner" width="100%" />
  
  <br />
  <br />

  <img src="assets/logo.png" alt="Interceptify Logo" width="150" />

  <h1>Interceptify</h1>

  <p>
    <strong>A Modern, Modular, and High-Performance MITM Framework</strong>
  </p>

  <p>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Language-Go_1.25+-00ADD8?style=flat&logo=go" alt="Go Version" /></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License" /></a>
    <a href="https://github.com/ismailtsdln/interceptify/actions"><img src="https://img.shields.io/badge/Build-Passing-success" alt="Build Status" /></a>
    <a href="https://goreportcard.com/report/github.com/ismailtsdln/interceptify"><img src="https://goreportcard.com/badge/github.com/ismailtsdln/interceptify" alt="Go Report Card" /></a>
  </p>
  
  <p>
    <em>Intercept, Analyze, and Manipulate Network Traffic with Ease.</em>
  </p>
</div>

---

## 📖 Overview

**Interceptify** is a next-generation Man-In-The-Middle (MITM) proxy tool developed in Go. Designed for security researchers, penetration testers, and developers, it offers a robust platform to inspect and modify HTTP/HTTPS and HTTP/2 traffic on the fly.

Unlike traditional proxies, Interceptify focuses on **modularity and performance**. Its plugin-based architecture allows you to write custom logic in Go to manipulate traffic without recompiling the core engine. With a sleek, real-time web dashboard, you can visualize interception usage like never before.

## 🚀 Key Features

- **🔒 HTTPS/TLS MITM**: Seamlessly intercept encrypted traffic with automatic, dynamic certificate generation.
- **⚡ HTTP/2 Support**: Native support for HTTP/2 multiplexing, ensuring modern web apps work flawlessly.
- **📊 Real-time Dashboard**: A stunning, glassmorphism-inspired web UI that monitors traffic in real-time using Server-Sent Events (SSE).
- **🧩 Modular Plugin System**: Extend functionality with simple Go plugins. Inject headers, drop packets, or modify payloads with just a few lines of code.
- **💻 Cross-Platform**: A single binary that runs on **macOS**, **Linux**, and **Windows**.
- **📝 Live Logging**: Insightful console logs and dashboard metrics keep you informed of every byte transferred.

## 🛠️ Installation

### Option 1: Using `go install` (Recommended)

The easiest way to install Interceptify is using the Go toolchain:

```bash
go install github.com/ismailtsdln/interceptify@latest
```

Ensure your `$(go env GOPATH)/bin` is in your system's `PATH`.

### Option 2: Build from Source

Requirements: Go 1.25+

```bash
git clone https://github.com/ismailtsdln/interceptify.git
cd interceptify
go build -o interceptify main.go
```

## 🚥 Quick Start

### 1. Start the Proxy

Launch Interceptify with default settings (listening on `127.0.0.1:8080`):

```bash
interceptify start
```

Or specify a custom address and port:

```bash
interceptify start --address 0.0.0.0 --port 9090
```

### 2. Configure Your Client

Set your browser or system proxy to `127.0.0.1:8080`.

### 3. Trust the CA Certificate

To intercept HTTPS traffic without warnings:
1.  Run Interceptify once to generate the CA.
2.  Locate the certificate at `~/.interceptify/ca.crt`.
3.  Import this file into your OS or browser's trusted root store.

### 4. Access the Dashboard

Open your web browser and navigate to:
👉 **[http://interceptify.local](http://interceptify.local)** (or `http://localhost:8080`)

## 🧩 Plugin Development

Interceptify is designed to be extensible. You can easily write plugins to manipulate traffic.

**Example: A Simple Header Injector**

```go
type MyPlugin struct {
    plugins.BasePlugin
}

func (p *MyPlugin) Name() string { return "HeaderInjector" }

func (p *MyPlugin) OnRequest(req *http.Request) (*http.Request, *http.Response) {
    req.Header.Set("X-Intercepted-By", "Interceptify")
    return req, nil
}
```

## 🗺️ Roadmap

- [x] Basic HTTP/HTTPS Interception
- [x] HTTP/2 Support
- [x] Web Dashboard (Real-time SSE)
- [x] `go install` Support
- [ ] WebSocket Interception
- [ ] Response Replay & Fuzzing
- [ ] Scriptable Python/Lua Plugins

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the repository
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## ⚖️ Legal Disclaimer

**Interceptify is for educational and authorized security testing purposes only.**
The authors are not responsible for any misuse or damage caused by this tool. tailored for ethical hacking and debugging.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/ismailtsdln">Ismail Tasdelen</a></p>
</div>
