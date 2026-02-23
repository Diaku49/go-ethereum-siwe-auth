package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	expiry time.Duration
	issuer string
}

func NewJWTService(secret string, expiry time.Duration, issuer string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
		expiry: expiry,
		issuer: issuer,
	}
}

type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

func (j *JWTService) Issue(address string) (string, error) {
	now := time.Now()

	claims := Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secret)
}

// This matches your middleware TokenVerifier interface: Verify(token)->address
func (j *JWTService) Verify(token string) (string, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return "", errors.New("invalid token")
	}

	// Optional issuer check
	if j.issuer != "" && claims.Issuer != j.issuer {
		return "", errors.New("invalid issuer")
	}

	if claims.Address == "" {
		return "", errors.New("missing address claim")
	}

	return claims.Address, nil
}
