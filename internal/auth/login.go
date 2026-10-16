package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrPassIncorect = errors.New("incorrect password")
)

type LoginInfo struct {
	Pass  string `json:"password,omitempty"`
	Token string `json:"token,omitempty"`
}

func Login(pswrd string) (string, error) {
	pass := os.Getenv("TODO_PASSWORD")
	if pswrd != pass {
		return "", ErrPassIncorect
	}

	passHash := sha256.Sum256([]byte(pass))

	secretKey := []byte(pass + "finalPrjсt")

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"hash": hex.EncodeToString(passHash[:])})

	signedToken, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
