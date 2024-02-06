package token

import (
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	panic("unimplemented")
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

func (payload *Payload) GetIssuer() (string, error) {
	return payload.IssuedAt.GoString(), nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

func (payload *Payload) GetSubject() (string, error) {
	panic("unimplemented")
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

func NewPayload(email string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}
