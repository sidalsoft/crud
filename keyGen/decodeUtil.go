package keyGen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
)

func HomeworkInit() {
	publicKeyBytes, err := os.ReadFile("public.key")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := decodePublicKey(publicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	planText := []byte("Go rullezzz!")
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, planText, nil)
	if err != nil {
		log.Fatal(err)
	}
	cipherTextFile, err := os.Create("ciphertext.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := cipherTextFile.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	_, _ = cipherTextFile.Write([]byte(fmt.Sprintf("%x", cipherText)))
}

func Init() {
	publicKeyBytes, err := os.ReadFile("public.key")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := decodePublicKey(publicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	planText := []byte("Go rullezzz!")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, planText, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%x", ciphertext)

	privateKeyByts, err := os.ReadFile("private.key")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := decodePrivateKey(privateKeyByts)
	if err != nil {
		log.Fatal(err)
	}
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", decryptedText)
}

func decodePublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("can't decode pem block")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func decodePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("can't decode pem block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
