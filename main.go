package main

import (
	"gitlab.com/transumption/unstable/seed-ca/ecdsa"

	"bytes"
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/argon2"
	"log"
	"math/big"
	"os"
	"time"
)

const salt = "jcD_:}VIy'_q)$>^@=q&2ywx6"

func main() {
	seed, ok := os.LookupEnv("SEED_CA")
	if !ok {
		log.Fatal(`
This program deterministically derives an ECDSA CA certificate and private key
from a seed (an arbitrary string). You need to set SEED_CA environment variable.

Please make sure to backup your seed. If you're using the seed to sign Android
apps, you won't be able to push any future updates if you ever lose it.`)
	}

	key := argon2.IDKey([]byte(seed), []byte(salt), 2, 1536*1024, 4, 73)
	rng := bytes.NewReader(key)

	priv, err := ecdsa.GenerateKey(elliptic.P521(), rng)
	if err != nil {
		log.Fatal(err)
	}

	privFile, err := os.Create("ca.key")
	if err != nil {
		log.Fatal(err)
	}

	defer privFile.Close()

	privX509, err := x509.MarshalECPrivateKey(priv.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = pem.Encode(privFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privX509})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(privFile.Name())

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			Organization:  []string{""},
			Country:       []string{""},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},

		NotBefore: time.Unix(0, 0),
		// https://tools.ietf.org/html/rfc5280#section-4.1.2.5
		NotAfter: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),

		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}

	cert, err := x509.CreateCertificate(rng, ca, ca, priv.Public(), priv)
	if err != nil {
		log.Fatal(err)
	}

	certFile, err := os.Create("ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(certFile.Name())
}
