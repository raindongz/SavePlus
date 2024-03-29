package token

import "time"

type Maker interface {
	//CreateToken create new token for specific username and duration
	CreateToken(id string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
