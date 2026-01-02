package JWT

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func LoadOrGenerateECDSAKeys() (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat("private_key.pem"); err == nil {
		return loadECDSAKeys()
	}
	return generateECDSAKeys()
}

func generateECDSAKeys() (*ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile("private_key.pem", pem.EncodeToMemory(&pem.Block{
		Type: "EC PRIVATE KEY",
		Bytes: privBytes,
	}), 0600); err != nil {
		return nil, err
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile("public_key.pem", pem.EncodeToMemory(&pem.Block{
		Type: "PUBLIC KEY",
		Bytes: pubBytes,
	}), 0644); err != nil {
		return nil, err
	}
	return priv, nil
}

func loadECDSAKeys() (*ecdsa.PrivateKey, error) {
	privPem, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("Couldn't load ecdsa private key %v\n", err)
	}
	block, _ := pem.Decode(privPem)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		log.Fatalf("Couldn't decode private key %v\n", err)
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Couldn't parse private key %v\n", err)
	}
	return priv, err
}