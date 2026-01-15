package plugins

import (
	"log"
	"net/http"
)

// Manager handles the lifecycle and execution of plugins
type Manager struct {
	plugins []Plugin
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make([]Plugin, 0),
	}
}

// Register adds a plugin to the manager
func (m *Manager) Register(p Plugin) {
	log.Printf("Registering plugin: %s", p.Name())
	m.plugins = append(m.plugins, p)
}

// RunRequestHooks runs all registered OnRequest hooks
func (m *Manager) RunRequestHooks(req *http.Request) (*http.Request, *http.Response) {
	for _, p := range m.plugins {
		var resp *http.Response
		req, resp = p.OnRequest(req)
		if resp != nil {
			// A plugin returned a response, short-circuit
			return req, resp
		}
	}
	return req, nil
}

// RunResponseHooks runs all registered OnResponse hooks
func (m *Manager) RunResponseHooks(req *http.Request, resp *http.Response) *http.Response {
	// Run in reverse order for response hooks? Usually yes.
	for i := len(m.plugins) - 1; i >= 0; i-- {
		resp = m.plugins[i].OnResponse(req, resp)
	}
	return resp
}
