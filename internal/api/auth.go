package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/moongo444/go_final_project/internal/auth"
	"github.com/moongo444/go_final_project/internal/db"
)

var (
	ErrPassEmpty = errors.New("password is empty")
	ErrJwtMethod = errors.New("incorrect JWT method")
	ErrJwtHash   = errors.New("JWT hash does not match")
)

type claims struct {
	Hash string
	jwt.RegisteredClaims
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	var buf bytes.Buffer
	var resp auth.LoginInfo

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	if len(resp.Pass) == 0 {
		writeJson(w, ErrPassEmpty.Error(), true, ErrPassEmpty)
		return
	}

	token, err := auth.Login(resp.Pass)
	if err != nil {
		writeJson(w, err.Error(), true, err)
		return
	}

	writeJson(w, token, true, nil)
}

func CheckAuth(next http.HandlerFunc, pass string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var token string
			var claim claims

			shaPass := sha256.Sum256([]byte(pass))

			hashPass := hex.EncodeToString(shaPass[:])

			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			token = cookie.Value

			if token == "" {
				http.Error(w, "token is empty", http.StatusUnauthorized)
				return
			}

			jwtToken, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (interface{}, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, ErrJwtMethod
				}
				return []byte(pass + "finalPrjсt"), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if jwtToken == nil {
				http.Error(w, "JWT is empty", http.StatusUnauthorized)
				return
			}

			if !jwtToken.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if claim.Hash != hashPass {
				http.Error(w, ErrJwtHash.Error(), http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
