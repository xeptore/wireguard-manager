package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func generateToken(secret []byte, userID, username, role string) (string, error) {
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
		Subject(userID).
		Claim("username", username).
		Claim("role", role).
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

func parseVerifyToken(secret []byte, token string) (string, error) {
	parseOptions := []jwt.ParseOption{
		jwt.WithVerify(jwa.HS512, secret),
		jwt.WithValidate(true),
		jwt.WithAudience("users"),
		jwt.WithIssuer("https://github.com/xeptore/wireguard-manager"),
		jwt.WithMinDelta(time.Second*10, jwt.ExpirationKey, jwt.IssuedAtKey),
	}
	parsedToken, err := jwt.Parse([]byte(token), parseOptions...)
	if nil != err {
		return "", fmt.Errorf("parsing and verifying token failed: %w", err)
	}

	return parsedToken.Subject(), nil
}
