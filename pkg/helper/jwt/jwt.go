package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	secret []byte
}

func NewJWTUtil(secret []byte) *JWTUtil {
	return &JWTUtil{secret: secret}
}

func (j *JWTUtil) CreateAccessToken(userID int64, tokenVersion int, remember bool) (string, time.Time, error) {
	tokenExpiration := 15 * time.Minute
	if remember {
		tokenExpiration = 7 * 24 * time.Hour
	}

	exp := time.Now().Add(tokenExpiration)

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           userID,
		"exp":           exp.Unix(),
		"token_version": tokenVersion,
	})

	accessToken, err := at.SignedString(j.secret)
	return accessToken, exp, err
}

func (j *JWTUtil) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRandomString() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
