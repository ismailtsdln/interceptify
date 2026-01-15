package ca

import (
	"os"
	"testing"
)

func TestCA(t *testing.T) {
	caCertPath := "test_ca.crt"
	caKeyPath := "test_ca.key"

	defer os.Remove(caCertPath)
	defer os.Remove(caKeyPath)

	// Test CA Creation/Generation
	caInstance, err := NewCA(caCertPath, caKeyPath)
	if err != nil {
		t.Fatalf("failed to create CA: %v", err)
	}

	if caInstance.Cert == nil || caInstance.Key == nil {
		t.Fatal("CA cert or key is nil")
	}

	// Test Certificate Signing
	host := "example.com"
	cert, key, err := caInstance.SignCertificate(host)
	if err != nil {
		t.Fatalf("failed to sign certificate for %s: %v", host, err)
	}

	if cert == nil || key == nil {
		t.Fatal("signed cert or key is nil")
	}

	if cert.Subject.CommonName != host {
		t.Errorf("expected CommonName %s, got %s", host, cert.Subject.CommonName)
	}
}
