package token

import "time"

type Maker interface {
	// create a new token for a user for a specific duration
	CreateToken(username string, duration time.Duration) (string, error)

	// check if the token is valid
	VerifyToken(token string) (*Payload, error)
}