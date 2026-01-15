package attack

import (
	"log"
	"net/http"

	"github.com/ismailtsdln/interceptify/pkg/plugins"
)

// LoggerPlugin is a built-in plugin for logging traffic
type LoggerPlugin struct {
	plugins.BasePlugin
}

func (p *LoggerPlugin) Name() string {
	return "Logger"
}

func (p *LoggerPlugin) Description() string {
	return "Logs all intercepted HTTP requests and responses"
}

func (p *LoggerPlugin) OnRequest(req *http.Request) (*http.Request, *http.Response) {
	log.Printf("[Plugin:Logger] Request: %s %s", req.Method, req.URL.String())
	return req, nil
}

func (p *LoggerPlugin) OnResponse(req *http.Request, resp *http.Response) *http.Response {
	log.Printf("[Plugin:Logger] Response: %d for %s", resp.StatusCode, req.URL.String())
	return resp
}
