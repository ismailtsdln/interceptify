package plugins

import (
	"net/http"
)

// Plugin defines the interface for Interceptify plugins
type Plugin interface {
	Name() string
	Description() string

	// Hooks
	OnRequest(req *http.Request) (*http.Request, *http.Response)
	OnResponse(req *http.Request, resp *http.Response) *http.Response
}

// BasePlugin provides a default implementation for the Plugin interface
type BasePlugin struct{}

func (p *BasePlugin) Name() string        { return "BasePlugin" }
func (p *BasePlugin) Description() string { return "Default base plugin" }
func (p *BasePlugin) OnRequest(req *http.Request) (*http.Request, *http.Response) {
	return req, nil
}
func (p *BasePlugin) OnResponse(req *http.Request, resp *http.Response) *http.Response {
	return resp
}
