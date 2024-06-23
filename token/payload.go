package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidHeader = errors.New("invalid header")
)

// These are claims
type Payload struct {
	ID        uuid.UUID        `json:"id"`
	Issuer    string           `json:"iss,omitempty"`
	Subject   string           `json:"sub,omitempty"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	Username  string           `json:"username"`
	Audience  jwt.ClaimStrings `json:"aud,omitempty"`
	IssuedAt  *jwt.NumericDate `json:"issued_at"`
	ExpiredAt *jwt.NumericDate `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiredAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}

	return payload, nil
}

// GetExpirationTime implements the Claims interface.
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return payload.ExpiredAt, nil
}

// GetIssuedAt implements the Claims interface.
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return payload.IssuedAt, nil
}

// GetNotBefore implements the Claims interface.
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return payload.NotBefore, nil
}

// GetAudience implements the Claims interface.
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return payload.Audience, nil
}

// GetIssuer implements the Claims interface.
func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

// GetSubject implements the Claims interface.
func (payload *Payload) GetSubject() (string, error) {
	return payload.Subject, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt.Time) {
		return ErrExpiredToken
	}

	return nil
}
