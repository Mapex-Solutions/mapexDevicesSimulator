package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// TLSMaterial is a self-contained PKI for the certificate-auth journey: one CA
// that signs both the broker's server cert and the device's client cert. The
// device trusts caPEM to verify the broker; the broker trusts the same CA to
// verify the device — classic mutual TLS, generated fresh per run so no fixture
// files live in the tree.
type TLSMaterial struct {
	CAPEM         string // CA cert; the device's tlsCaPem and the broker's client-CA pool
	ServerCertPEM string // broker server cert (SAN 127.0.0.1)
	ServerKeyPEM  string // broker server key
	ClientCertPEM string // device client cert
	ClientKeyPEM  string // device client key
}

// GenerateTLSMaterial mints the CA and the server/client leaf certs for the
// mutual-TLS MQTT journey. Certs are valid for localhost (127.0.0.1) so the
// in-process broker and the sidecar — both on the loopback — verify cleanly.
func GenerateTLSMaterial() (TLSMaterial, error) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return TLSMaterial{}, fmt.Errorf("ca key: %w", err)
	}
	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "e2e-mqtt-ca"},
		NotBefore:             time.Unix(0, 0),
		NotAfter:              time.Unix(1<<31-1, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	if err != nil {
		return TLSMaterial{}, fmt.Errorf("ca cert: %w", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		return TLSMaterial{}, fmt.Errorf("parse ca: %w", err)
	}

	serverCertPEM, serverKeyPEM, err := leaf(caCert, caKey, "e2e-mqtt-broker", x509.ExtKeyUsageServerAuth, true)
	if err != nil {
		return TLSMaterial{}, fmt.Errorf("server leaf: %w", err)
	}
	clientCertPEM, clientKeyPEM, err := leaf(caCert, caKey, "e2e-mqtt-device", x509.ExtKeyUsageClientAuth, false)
	if err != nil {
		return TLSMaterial{}, fmt.Errorf("client leaf: %w", err)
	}

	return TLSMaterial{
		CAPEM:         pemBlock("CERTIFICATE", caDER),
		ServerCertPEM: serverCertPEM,
		ServerKeyPEM:  serverKeyPEM,
		ClientCertPEM: clientCertPEM,
		ClientKeyPEM:  clientKeyPEM,
	}, nil
}

// leaf mints a single leaf cert signed by the CA. withLoopbackSAN adds the
// 127.0.0.1 / localhost SANs the server cert needs for the sidecar to verify
// the broker's identity; the client cert is authenticated by signature alone.
func leaf(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, cn string, eku x509.ExtKeyUsage, withLoopbackSAN bool) (certPEM, keyPEM string, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("key: %w", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Unix(0, 0).UnixNano() + int64(len(cn))),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Unix(0, 0),
		NotAfter:     time.Unix(1<<31-1, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{eku},
	}
	if withLoopbackSAN {
		tmpl.IPAddresses = []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
		tmpl.DNSNames = []string{"localhost"}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
	if err != nil {
		return "", "", fmt.Errorf("cert: %w", err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", fmt.Errorf("marshal key: %w", err)
	}
	return pemBlock("CERTIFICATE", der), pemBlock("PRIVATE KEY", keyDER), nil
}

// pemBlock encodes DER bytes into a PEM string of the given block type.
func pemBlock(typ string, der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}))
}
