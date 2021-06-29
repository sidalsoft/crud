package keyGen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func GenCertificate() {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"Go alif Academy"},
			Country:      []string{"TJ"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(3, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	reader := rand.Reader
	bitSize := 4096

	caKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Print(err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		log.Print(err)
	}

	err = encodePrivateKey(caKey, "ca-private.key")
	if err != nil {
		log.Print(err)
	}

	err = encodePublicKey(&caKey.PublicKey, "ca-publick.key")
	if err != nil {
		log.Print(err)
	}

	err = encodeCert(caBytes, "ca.crt")
	if err != nil {
		log.Print(err)
	}

	//for server
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:       []string{"Go Alif Academy"},
			OrganizationalUnit: []string{"Dev"},
			Country:            []string{"TJ"},
		},
		DNSNames:    []string{"go.alif.hack."},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		//SubjectKeyId:[]byte{1,2,3,4,6},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Print(err)
	}

	//Издаем сертификат для сервераб подписанный приватным ключом CА и указываем родительский сертификат

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certKey.PublicKey, caKey)
	if err != nil {
		log.Print(err)
	}

	err = encodePrivateKey(certKey, "server-private.key")
	if err != nil {
		log.Print(err)
	}

	err = encodePublicKey(&certKey.PublicKey, "server-public.key")
	if err != nil {
		log.Print(err)
	}

	err = encodeCert(certBytes, "server.crt")
	if err != nil {
		log.Print(err)
	}

}

func encodeCert(cert []byte, path string) error {
	certFile, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := certFile.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	data := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	err = pem.Encode(certFile, data)
	if err != nil {
		return err
	}
	return nil
}
