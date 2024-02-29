package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type GenerateParams struct {
	SigningKey string
	TTL        time.Duration
	Subject    string
	Issuer     string
}

func Generate(params GenerateParams) (string, error) {
	if params.TTL < time.Minute {
		return "", errors.New("TTL must be at least 1 minute")
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   params.Subject,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(params.TTL)),
		Issuer:    params.Issuer,
		Audience:  jwt.ClaimStrings{params.Issuer},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(params.SigningKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type ParseParams struct {
	SigningKey string
	Issuer     string
}

func Parse(signedToken string, params ParseParams) (subject string, err error) {
	opts := []jwt.ParserOption{
		jwt.WithIssuer(params.Issuer), jwt.WithAudience(params.Issuer),
		jwt.WithLeeway(time.Minute), jwt.WithIssuedAt(), jwt.WithExpirationRequired(),
	}

	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(params.SigningKey), nil
	}, opts...)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", ErrInvalidToken
}
