package jwtkit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

type Encryptor interface {
	generatingPublicPrivateToPEM() error
	gettingPublicFromPEM(string) (*ecdsa.PublicKey, error)
	gettingPrivateFromPEM(string) (*ecdsa.PrivateKey, error)
}

func GeneratePublicPrivateToPEM(e Encryptor) error {
	err := e.generatingPublicPrivateToPEM()
	if err != nil {
		return err
	}
	return nil
}

func GetPublicFromPEM(e Encryptor) (*ecdsa.PublicKey, error) {
	enc, ok := e.(*ECDSA)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	publicKey, err := e.gettingPublicFromPEM(enc.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func GetPrivateFromPEM(e Encryptor) (*ecdsa.PrivateKey, error) {
	enc, ok := e.(*ECDSA)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	privateKey, err := e.gettingPrivateFromPEM(enc.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (enc *ECDSA) generatingPublicPrivateToPEM() error {
	_, errPrivate := os.Stat(enc.PrivateKeyPath)
	_, errPublic := os.Stat(enc.PublicKeyPath)

	if os.IsNotExist(errPrivate) && os.IsNotExist(errPublic) {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}

		publicKey := &privateKey.PublicKey

		pemPrivateKeyFile, err := os.Create(enc.PrivateKeyPath)
		if err != nil {
			return err
		}

		pemPublicKeyFile, err := os.Create(enc.PublicKeyPath)
		if err != nil {
			return err
		}

		defer pemPrivateKeyFile.Close()
		defer pemPublicKeyFile.Close()

		marshalPrivate, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return err
		}

		pemPrivateKeyBlock := &pem.Block{
			Type:  "E256 PRIVATE KEY",
			Bytes: marshalPrivate,
		}

		asn1Bytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return err
		}

		pemPublicKeyBlock := &pem.Block{
			Type:  "E256 PUBLIC KEY",
			Bytes: asn1Bytes,
		}

		err = pem.Encode(pemPrivateKeyFile, pemPrivateKeyBlock)
		if err != nil {
			return err
		}

		err = pem.Encode(pemPublicKeyFile, pemPublicKeyBlock)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (enc *ECDSA) gettingPublicFromPEM(string) (*ecdsa.PublicKey, error) {
	publicKeyData, err := ioutil.ReadFile(enc.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	publicPemBlock, _ := pem.Decode(publicKeyData)
	if publicPemBlock == nil || publicPemBlock.Type != "E256 PUBLIC KEY" {
		return nil, errors.New("Failed to decode PEM block containing public key")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicPemBlock.Bytes)
	if err != nil {
		if cert, err := x509.ParseCertificate(publicPemBlock.Bytes); err == nil {
			parsedPublicKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	publicKey, ok := parsedPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return publicKey, nil
}

func (enc *ECDSA) gettingPrivateFromPEM(string) (*ecdsa.PrivateKey, error) {
	privateKeyData, err := ioutil.ReadFile(enc.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	privatePemBlock, _ := pem.Decode(privateKeyData)
	if privatePemBlock == nil || privatePemBlock.Type != "E256 PRIVATE KEY" {
		return nil, errors.New("Failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(privatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
