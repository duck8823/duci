package application

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/tcnksm/go-latest"
	"os"
	"testing"
)

type MaskString = maskString

func SetVersion(ver string) (reset func()) {
	tmp := version
	version = ver
	return func() {
		version = tmp
	}
}

func SetCheckResponse(chr *latest.CheckResponse) (reset func()) {
	tmp := checked
	checked = chr
	return func() {
		checked = tmp
	}
}

func CheckLatestVersion() {
	checkLatestVersion()
}

func TrimSuffix(tag string) string {
	return trimSuffix(tag)
}

func GetCtxKey() *string {
	return &ctxKey
}

func GenerateSSHKey(t *testing.T, path string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		t.Fatalf("error occur: %+v", err)
	}

	if err := key.Validate(); err != nil {
		t.Fatalf("error occur: %+v", err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		t.Fatalf("error occur: %+v", err)
	}

	if _, err := f.Write(pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	})); err != nil {
		t.Fatalf("error occur: %+v", err)
	}
}
