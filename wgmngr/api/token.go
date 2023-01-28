package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type TokenClaims struct {
	UserID, Username, Role string
}

func generateToken(secret []byte, claims TokenClaims) (string, error) {
	tokenID, err := uuid.NewV4()
	if nil != err {
		return "", fmt.Errorf("unable to generate token id: %w", err)
	}

	token, err := jwt.NewBuilder().
		Issuer("https://github.com/xeptore/wireguard-manager").
		IssuedAt(time.Now()).
		Audience([]string{"users"}).
		Expiration(time.Now().Add(time.Minute*15)).
		JwtID(tokenID.String()).
		NotBefore(time.Now()).
		Subject(claims.UserID).
		Claim("username", claims.Username).
		Claim("role", claims.Role).
		Build()
	if nil != err {
		return "", fmt.Errorf("failed to build token: %w", err)
	}

	opts := []jwt.SignOption{}
	serialized, err := jwt.Sign(token, jwa.HS512, []byte(secret), opts...)
	if nil != err {
		return "", fmt.Errorf("unable to sign auth token: %w", err)
	}

	var b strings.Builder
	b.Grow(len(serialized))
	if _, err := b.Write(serialized); nil != err {
		return "", err
	}

	return b.String(), nil
}

func parseVerifyToken(secret []byte, token string) (*TokenClaims, error) {
	parseOptions := []jwt.ParseOption{
		jwt.WithVerify(jwa.HS512, secret),
		jwt.WithValidate(true),
		jwt.WithAudience("users"),
		jwt.WithIssuer("https://github.com/xeptore/wireguard-manager"),
		jwt.WithMinDelta(time.Second*10, jwt.ExpirationKey, jwt.IssuedAtKey),
	}
	parsedToken, err := jwt.Parse([]byte(token), parseOptions...)
	if nil != err {
		return nil, fmt.Errorf("parsing and verifying token failed: %w", err)
	}

	roleClaim, ok := parsedToken.Get("role")
	if !ok {
		return nil, errors.New("role claim expected to exist")
	}
	role, ok := roleClaim.(string)
	if !ok {
		return nil, errors.New("username is expected to be a string")
	}

	usernameClaim, ok := parsedToken.Get("username")
	if !ok {
		return nil, errors.New("role claim expected to exist")
	}
	username, ok := usernameClaim.(string)
	if !ok {
		return nil, errors.New("username is expected to be a string")
	}

	return &TokenClaims{
		UserID:   parsedToken.Subject(),
		Role:     role,
		Username: username,
	}, nil
}
