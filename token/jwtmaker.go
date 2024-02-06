package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	symmetricKey []byte
}

func (maker *JWTMaker) CreateToken(email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString(maker.symmetricKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return maker.symmetricKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &JWTMaker{symmetricKey: []byte(secretKey)}, nil
}
