package attack

import (
	"log"
	"net/http"

	"github.com/ismailtsdln/interceptify/pkg/parser"
	"github.com/ismailtsdln/interceptify/pkg/plugins"
)

// ModifierPlugin is a built-in plugin for modifying traffic
type ModifierPlugin struct {
	plugins.BasePlugin
	manipulator *parser.Manipulator
}

func (p *ModifierPlugin) Name() string {
	return "Modifier"
}

func (p *ModifierPlugin) Description() string {
	return "Injects headers and modifies response bodies"
}

func (p *ModifierPlugin) OnRequest(req *http.Request) (*http.Request, *http.Response) {
	log.Printf("[Plugin:Modifier] Injecting X-Interceptified header")
	p.manipulator.InjectHeader(req.Header, "X-Interceptified", "true")
	return req, nil
}

func (p *ModifierPlugin) OnResponse(req *http.Request, resp *http.Response) *http.Response {
	log.Printf("[Plugin:Modifier] Checking for body modification")
	// Example: Replace "Google" with "Interceptify" in response bodies
	err := p.manipulator.ReplaceInBody(resp, "Google", "Interceptify")
	if err != nil {
		log.Printf("[Plugin:Modifier] Body replacement failed: %v", err)
	}
	return resp
}

// NewModifierPlugin creates a new ModifierPlugin
func NewModifierPlugin() *ModifierPlugin {
	return &ModifierPlugin{
		manipulator: parser.NewManipulator(),
	}
}
