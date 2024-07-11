package token

import "time"

// interface for tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	ValidateToken(token string) (*Payload, error)
}
