package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/smallstep/pkcs7"
)

func GenerateKeys(certPath, keyPath string) error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"My Sign Service"}, // ?
			CommonName:   "User Name",                 // ?
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return err
	}

	// Сохраняем ключ
	keyFile, _ := os.Create(keyPath)
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

	// Сохраняем сертификат
	certFile, _ := os.Create(certPath)
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	return nil
}

// SignDocument создает отсоединенную PKCS7 подпись
func SignDocument(data []byte, certPath, keyPath string) ([]byte, error) {
	// Загружаем сертификат
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certData)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Загружаем ключ
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyData)
	privKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Подписываем
	signedData, err := pkcs7.NewSignedData(data)
	if err != nil {
		return nil, err
	}

	if err := signedData.AddSigner(cert, privKey, pkcs7.SignerInfoConfig{}); err != nil {
		return nil, err
	}

	signedData.Detach()
	return signedData.Finish()
}

// VerifySignature проверяет, соответствует ли подпись файлу
func VerifySignature(data []byte, signature []byte) error {
	p7, err := pkcs7.Parse(signature)
	if err != nil {
		return err
	}
	p7.Content = data // Приаттачиваем файл обратно для проверки хэша
	return p7.Verify()
}
