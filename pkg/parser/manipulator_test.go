package parser

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestReplaceInBody(t *testing.T) {
	m := NewManipulator()

	originalBody := "Hello World"
	resp := &http.Response{
		Body:   io.NopCloser(bytes.NewBufferString(originalBody)),
		Header: make(http.Header),
	}
	resp.Header.Set("Content-Length", "11")

	err := m.ReplaceInBody(resp, "World", "Interceptify")
	if err != nil {
		t.Fatalf("ReplaceInBody failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	expectedBody := "Hello Interceptify"
	if string(body) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, string(body))
	}

	expectedLen := "18"
	if resp.Header.Get("Content-Length") != expectedLen {
		t.Errorf("expected Content-Length %q, got %q", expectedLen, resp.Header.Get("Content-Length"))
	}
}
