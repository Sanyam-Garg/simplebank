package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil{
		return "", err 
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(jwtMaker.secretKey))
}

func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// Verify that the signing algorithm matches with the oncoming header
	keyFunc := func (token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, ErrInvalidToken
		}

		return []byte(jwtMaker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil{
		verr, ok := err.(*jwt.ValidationError)
		// Check if it is the errExpiredToken thrown by the Valid function in payload
		if ok && errors.Is(verr.Inner, ErrExpiredToken){
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok{
		return nil, ErrInvalidToken
	}

	return payload, nil
}
