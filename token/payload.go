package token

import (
	"fmt"
	"time"

	 "github.com/google/uuid"
)

type Payload struct {
	ID uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time    `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error){
	tokenID, err := uuid.NewRandom()
	if err != nil{
		return nil, err
	}

	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return fmt.Errorf("token has expired")
	}
	return nil
}
