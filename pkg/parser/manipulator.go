package parser

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// Manipulator provides utility functions for HTTP packet manipulation
type Manipulator struct{}

// NewManipulator creates a new Manipulator instance
func NewManipulator() *Manipulator {
	return &Manipulator{}
}

// ReplaceInBody replaces all occurrences of 'old' with 'new' in the response body
func (m *Manipulator) ReplaceInBody(resp *http.Response, old, new string) error {
	if resp.Body == nil {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	newBody := strings.ReplaceAll(string(body), old, new)
	resp.Body = io.NopCloser(bytes.NewBufferString(newBody))
	resp.ContentLength = int64(len(newBody))

	// Update Content-Length header if it exists
	if resp.Header.Get("Content-Length") != "" {
		resp.Header.Set("Content-Length", string(rune(len(newBody))))
	}

	return nil
}

// InjectHeader adds a header to the request or response
func (m *Manipulator) InjectHeader(header http.Header, key, value string) {
	header.Set(key, value)
}

// DropHeader removes a header from the request or response
func (m *Manipulator) DropHeader(header http.Header, key string) {
	header.Del(key)
}
