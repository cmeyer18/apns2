package certificate

import (
	"crypto/tls"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PKCS#12

func TestValidCertificateFromP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid.p12", "")
	assert.Nil(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestValidCertificateFromP12Bytes(t *testing.T) {
	bytes, _ := os.ReadFile("_fixtures/certificate-valid.p12")
	cer, err := FromP12Bytes(bytes, "")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestEncryptedValidCertificateFromP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid-encrypted.p12", "password")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestNoSuchFileP12File(t *testing.T) {
	cer, err := FromP12File("", "")
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
	assert.Equal(t, tls.Certificate{}, cer)
}

func TestBadPasswordP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid-encrypted.p12", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("pkcs12: decryption password incorrect").Error(), err.Error())
}

// PEM

func TestValidCertificateFromPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid.pem", "")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestValidCertificateFromPemBytes(t *testing.T) {
	bytes, _ := os.ReadFile("_fixtures/certificate-valid.pem")
	cer, err := FromPemBytes(bytes, "")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestValidCertificateFromPemFileWithPKCS8PrivateKey(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid-pkcs8.pem", "")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestValidCertificateFromPemBytesWithPKCS8PrivateKey(t *testing.T) {
	bytes, _ := os.ReadFile("_fixtures/certificate-valid-pkcs8.pem")
	cer, err := FromPemBytes(bytes, "")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestEncryptedValidCertificateFromPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid-encrypted.pem", "password")
	assert.NoError(t, err)
	assert.NotEqual(t, tls.Certificate{}, cer)
}

func TestNoSuchFilePemFile(t *testing.T) {
	cer, err := FromPemFile("", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
}

func TestBadPasswordPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid-encrypted.pem", "badpassword")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrFailedToDecryptKey, err)
}

func TestBadKeyPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-bad-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrFailedToParsePrivateKey, err)
}

func TestNoKeyPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-no-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrNoPrivateKey, err)
}

func TestNoCertificatePemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-no-pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrNoCertificate, err)
}
