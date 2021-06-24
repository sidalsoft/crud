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
	err = encodePrivateKey(err, key)
	if err != nil {
		log.Fatal(err)
	}
	err = encodePublicKey(err, &key.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
}

func encodePrivateKey(err error, key *rsa.PrivateKey) error {
	privateKeyFile, err := os.Create("private.key")
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

func encodePublicKey(err error, key *rsa.PublicKey) error {
	publicKeyFile, err := os.Create("public.key")
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
