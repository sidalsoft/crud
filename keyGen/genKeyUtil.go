package keyGen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"log"
	"os"
)

func KeyGenInit() {
	reader := rand.Reader
	butSize := 4096

	key, err := rsa.GenerateKey(reader, butSize)
	if err != nil {
		log.Fatal(err)
	}
	err = encodePrivateKey(key, "private.key")
	if err != nil {
		log.Fatal(err)
	}
	err = encodePublicKey(&key.PublicKey, "public.key")
	if err != nil {
		log.Fatal(err)
	}
}

func encodePrivateKey(key *rsa.PrivateKey, path string) error {
	privateKeyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := privateKeyFile.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	privateKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func encodePublicKey(key *rsa.PublicKey, path string) error {
	publicKeyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := publicKeyFile.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	asn1Bytes, err := asn1.Marshal(*key)
	publicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	err = pem.Encode(publicKeyFile, publicKey)
	if err != nil {
		return err
	}
	return nil
}
