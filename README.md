# Interceptify 📈

Interceptify is a modern, modular, and high-performance Man-In-The-Middle (MITM) framework developed in Go. It allows security researchers and developers to intercept, analyze, and manipulate network traffic (HTTP/HTTPS) with ease.

## 🚀 Features

- **HTTPS/TLS MITM**: Seamlessly intercept encrypted traffic using dynamic certificate generation.
- **Modular Plugin System**: Extend functionality with custom hooks for requests and responses.
- **Cross-Platform**: Built in Go, it runs on Linux, macOS, and Windows.
- **Built-in Logging**: Real-time traffic logging with the default logger plugin.
- **High Performance**: Leverages Go's concurrency model for efficient traffic handling.

## 🛠️ Components

- **CA Management**: Automatically generates a root CA and signs certificates for any host on-the-fly.
- **Proxy Engine**: A robust TCP listener that handles connection hijacking and TLS termination.
- **Plugin Manager**: A flexible architecture to register and run traffic manipulation logic.

## 📦 Installation

To build Interceptify from source, you need Go 1.21 or higher installed.

```bash
git clone https://github.com/ismailtsdln/interceptify.git
cd interceptify
go build -o interceptify main.go
```

## 🚥 Usage

### Start the Proxy

Launch Interceptify with default settings (listens on `127.0.0.1:8080`):

```bash
./interceptify start
```

You can specify a custom port and address:

```bash
./interceptify start --port 9090 --address 0.0.0.0
```

### Trust the Root CA

To intercept HTTPS traffic, you must trust the generated Root CA:

1. Locate the CA certificate at `~/.interceptify/ca.crt`.
2. Install it in your system or browser's trusted certificate store.

## 🧩 Plugin Development

Creating a plugin is simple. Implement the `Plugin` interface:

```go
type MyPlugin struct {
    plugins.BasePlugin
}

func (p *MyPlugin) Name() string { return "MyPlugin" }

func (p *MyPlugin) OnRequest(req *http.Request) (*http.Request, *http.Response) {
    // Modify request here
    return req, nil
}
```

## ⚖️ Legal Disclaimer

Interceptify is intended for educational purposes and authorized security testing only. Using this tool to intercept traffic without permission is illegal and unethical. The authors are not responsible for any misuse.

## 📄 License

This project is licensed under the MIT License.
