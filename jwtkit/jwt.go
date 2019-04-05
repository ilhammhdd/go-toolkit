package jwtkit

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"time"
)

type JWTString string

type JWT struct {
	Header    *Header    `json:"header"`
	Payload   *Payload   `json:"payload"`
	Signature *Signature `json:"signature"`
}

type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type Payload struct {
	Id        string                 `json:"jti,omitempty"`
	Audience  string                 `json:"aud,omitempty"`
	Issuer    string                 `json:"iss,omitempty"`
	IssuedAt  int64                  `json:"iat,omitempty"`
	ExpiredAt int64                  `json:"exp,omitempty"`
	NotBefore int64                  `json:"nbf,omitempty"`
	Claims    map[string]interface{} `json:"claims,omitempty"`
}

type Signature struct {
	Hashed []byte   `json:"hashed"`
	R      *big.Int `json:"r"`
	S      *big.Int `json:"s"`
}

func (je JWTExpiration) GenerateSignedJWTString(encrypt *ECDSA, audience string, issuer string, claims ...*map[string]interface{}) (JWTString, error) {
	privateKey, err := GetPrivateFromPEM(encrypt)
	if err != nil {
		return "", err
	}
	publicKey, err := GetPublicFromPEM(encrypt)
	if err != nil {
		return "", err
	}

	jwt := &JWT{
		Header: &Header{
			Algorithm: "ECDSA",
			Type:      "JWT",
		},
		Payload: &Payload{
			Audience:  audience,
			Issuer:    issuer,
			IssuedAt:  time.Now().UnixNano() / 1000000,
			ExpiredAt: (time.Now().UnixNano() / 1000000) + int64(je),
		},
	}

	if len(claims) != 0 {
		jwt.Payload.Claims = *claims[0]
	}

	jsonHeader, err := json.Marshal(jwt.Header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(jsonHeader)

	jsonPayload, err := json.Marshal(jwt.Payload)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(jsonPayload)

	hashBase := encodedHeader + "." + encodedPayload

	hasher := sha256.New()
	hasher.Write([]byte(hashBase))
	hashed := hasher.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return "", err
	}

	if !ecdsa.Verify(publicKey, hashed, r, s) {
		return "", errors.New("Generated hash, r, and s fails to verify")
	}

	jwt.Signature = &Signature{}
	jwt.Signature.Hashed = hashed
	jwt.Signature.R = r
	jwt.Signature.S = s
	jsonSignature, err := json.Marshal(jwt.Signature)
	if err != nil {
		return "", err
	}
	encodedSignature := base64.URLEncoding.EncodeToString(jsonSignature)

	signedJWT := JWTString(strings.Join([]string{encodedHeader, encodedPayload, encodedSignature}, "."))

	return signedJWT, nil
}

func hashedFromJWT(jwt *JWT) ([]byte, error) {

	jsonHeader, err := json.Marshal(jwt.Header)
	if err != nil {
		return nil, err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(jsonHeader)

	jsonPayload, err := json.Marshal(jwt.Payload)
	if err != nil {
		return nil, err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(jsonPayload)

	hashBase := encodedHeader + "." + encodedPayload

	hasher := sha256.New()
	hasher.Write([]byte(hashBase))
	hashed := hasher.Sum(nil)

	return hashed, nil
}

func VerifyJWTString(encrypt *ECDSA, j JWTString) (bool, error) {

	publicKey, err := GetPublicFromPEM(encrypt)
	if err != nil {
		return false, err
	}

	jwtParts := strings.Split(string(j), ".")

	hashBase := jwtParts[0] + "." + jwtParts[1]

	hasher := sha256.New()
	hasher.Write([]byte(hashBase))
	supposedHashed := hasher.Sum(nil)

	decodedSignature, err := base64.URLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return false, err
	}

	var signature Signature
	err = json.Unmarshal([]byte(decodedSignature), &signature)
	if err != nil {
		return false, err
	}

	if string(supposedHashed) != string(signature.Hashed) {
		return false, errors.New("invalid signature")
	}

	if !ecdsa.Verify(publicKey, signature.Hashed, signature.R, signature.S) {
		return false, err
	}

	signature.Hashed[0] ^= 0xff
	if ecdsa.Verify(publicKey, signature.Hashed, signature.R, signature.S) {
		return false, errors.New("Verify always true")
	}

	return true, nil
}

func ValidateExpired(j JWTString) (bool, error) {
	jwtParts := strings.Split(string(j), ".")
	decodedPayload, err := base64.URLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return false, err
	}

	var payload Payload
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return false, err
	}

	if now := time.Now().UnixNano() / 1000000; payload.ExpiredAt <= now {

		return false, err
	}

	return true, err
}

func RegenerateToken(encrypt *ECDSA, audience string, issuer string, je JWTExpiration, j JWTString) (JWTString, error) {
	jwt, err := GetJWT(j)
	if err != nil {
		return "", err
	}
	return je.GenerateSignedJWTString(encrypt, audience, issuer, &jwt.Payload.Claims)
}

func GetJWT(j JWTString) (*JWT, error) {
	jwtParts := strings.Split(string(j), ".")

	decodedHeader, err := base64.URLEncoding.DecodeString(jwtParts[0])
	if err != nil {
		return nil, err
	}
	decodedPayload, err := base64.URLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return nil, err
	}
	decodedSignature, err := base64.URLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return nil, err
	}

	var jwtHeader Header
	err = json.Unmarshal(decodedHeader, &jwtHeader)
	if err != nil {
		return nil, err
	}

	var jwtPayload Payload
	err = json.Unmarshal(decodedPayload, &jwtPayload)
	if err != nil {
		return nil, err
	}

	var jwtSignature Signature
	err = json.Unmarshal(decodedSignature, &jwtSignature)
	if err != nil {
		return nil, err
	}

	return &JWT{
		Header:    &jwtHeader,
		Payload:   &jwtPayload,
		Signature: &jwtSignature,
	}, nil
}
