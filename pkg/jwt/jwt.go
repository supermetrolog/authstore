package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const DEFAULT_SECRET_KEY = "kFkdskpodf33ldFDFd9fDfdLflmvz"

type Header struct {
	Alg string `json:"alg"` //alghoritm
	Typ string `json:"typ"` //token type
}
type Payload struct {
	Sub    string `json:"sub"`     // user id
	Iss    string `json:"iss"`     // app id
	Exp    int64  `json:"exp"`     // expire
	Jti    string `json:"jti"`     // JWT id
	Iat    int64  `json:"iat"`     // createdAt
	UserID int64  `json:"user_id"` // user id
}
type Token struct {
	Header    Header  `json:"header"`
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
	Secret    string  `json:"-"`
}

func New(h Header, p Payload, secret string) *Token {
	if h.Alg == "" {
		h.Alg = "HS256"
	}
	h.Typ = "JWT"
	if secret == "" {
		secret = DEFAULT_SECRET_KEY
	}
	return &Token{
		Header:  h,
		Payload: p,
		Secret:  secret,
	}
}

func (t *Token) Generate() (string, error) {
	headerBytes, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	payloadBytes, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}
	base64Header := base64.StdEncoding.EncodeToString(headerBytes)
	base64Payload := base64.StdEncoding.EncodeToString(payloadBytes)

	unsignedToken := base64Header + "." + base64Payload

	h := hmac.New(sha256.New, []byte(t.Secret))
	_, err = h.Write([]byte(unsignedToken))
	if err != nil {
		return "", err
	}
	sigSum := h.Sum(nil)
	fmt.Println("Secret:", t.Secret, len(t.Secret))
	fmt.Println("Sig sum: ", sigSum)
	signature := hex.EncodeToString(sigSum)

	fmt.Printf("Encoded sig sum %s\n", signature)
	token := base64Header + "." + base64Payload + "." + signature
	return token, nil
}
