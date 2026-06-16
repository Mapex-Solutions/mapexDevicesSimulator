package session

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

// selfSignedPEM mints a throwaway self-signed cert/key pair in PEM, so the test
// exercises buildTLSConfig against material with the same shape a real device
// would carry, without reaching for any fixture files.
func selfSignedPEM(t *testing.T) (certPem, keyPem string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test-device"},
		NotBefore:    time.Unix(0, 0),
		NotAfter:     time.Unix(1<<31-1, 0),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}
	keyDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	certPem = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	keyPem = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDer}))
	return certPem, keyPem
}

func TestBuildTLSConfig(t *testing.T) {
	certPem, keyPem := selfSignedPEM(t)

	tests := []struct {
		name        string
		cert        string
		key         string
		ca          string
		wantErr     bool
		wantClient  bool
		wantRootCAs bool
	}{
		{
			name:       "client keypair only",
			cert:       certPem,
			key:        keyPem,
			wantClient: true,
		},
		{
			name:        "ca only",
			ca:          certPem,
			wantRootCAs: true,
		},
		{
			name:        "client keypair and ca",
			cert:        certPem,
			key:         keyPem,
			ca:          certPem,
			wantClient:  true,
			wantRootCAs: true,
		},
		{
			name:    "broken keypair",
			cert:    "not a pem",
			key:     "not a pem",
			wantErr: true,
		},
		{
			name:    "invalid ca pem",
			ca:      "not a pem",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := buildTLSConfig(tt.cert, tt.key, tt.ca)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.MinVersion == 0 {
				t.Errorf("expected a minimum TLS version to be set")
			}
			if got := len(cfg.Certificates) == 1; got != tt.wantClient {
				t.Errorf("client cert present = %v, want %v", got, tt.wantClient)
			}
			if got := cfg.RootCAs != nil; got != tt.wantRootCAs {
				t.Errorf("root CAs present = %v, want %v", got, tt.wantRootCAs)
			}
		})
	}
}
