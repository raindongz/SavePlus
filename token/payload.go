package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// different types of error returned by the VerifiyTOken function
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains payload data for token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Uid       string    `json:"uid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with user id and duration
func NewPayload(uid string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Uid:       uid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	// fmt.Println("hello time is :", payload.ExpiredAt)
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
