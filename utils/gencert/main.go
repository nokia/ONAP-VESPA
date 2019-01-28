package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"
)

var (
	caFile   = flag.String("cafile", "ca.pem", "Path to CA certificate file")
	caKey    = flag.String("cakey", "cakey.pem", "Path to CA certificate's key file")
	certFile = flag.String("cert", "cert.pem", "Path to certificate file")
	keyFile  = flag.String("key", "key.pem", "Path to certificate's key file")
	hosts    = flag.String("hosts", "", "Coma separated list of hosts to insert in certificate")
)

func generateCertificate(certFile, keyFile string, hosts []string, parent *x509.Certificate, pKey *rsa.PrivateKey, isCA bool) (*x509.Certificate, *rsa.PrivateKey, error) {
	fmt.Println("Generating Key Pair")
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Nokia"},
		},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
	}
	template.SerialNumber, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	fmt.Println("Generating certificate")
	if parent == nil {
		parent = &template
		pKey = key
	}
	certDer, err := x509.CreateCertificate(rand.Reader, &template, parent, &key.PublicKey, pKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certDer)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Saving certificate")
	data := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	if err := ioutil.WriteFile(certFile, data, 0); err != nil {
		return nil, nil, err
	}
	fmt.Println("Saving key")
	data = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err := ioutil.WriteFile(keyFile, data, 0); err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

func generateCA(caFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	return generateCertificate(caFile, keyFile, nil, nil, nil, true)
}

func main() {
	flag.Parse()
	caCert, caKey, err := generateCA(*caFile, *caKey)
	if err != nil {
		panic(err)
	}
	hostList := strings.Split(*hosts, ",")
	if _, _, err := generateCertificate(*certFile, *keyFile, hostList, caCert, caKey, false); err != nil {
		panic(err)
	}
}
