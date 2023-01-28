package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func generateToken(secret []byte, username string) (string, error) {
	tokenID, err := uuid.NewV4()
	if nil != err {
		return "", fmt.Errorf("unable to generate token id: %w", err)
	}

	token := jwt.New()
	if err := token.Set(jwt.IssuerKey, "https://github.com/xeptore/wireguard-manager"); nil != err {
		return "", fmt.Errorf("unable to set auth token issuer key: %w", err)
	}

	if err := token.Set(jwt.AudienceKey, []string{"users"}); nil != err {
		return "", fmt.Errorf("unable to set auth token audience key: %w", err)
	}

	exp := time.Now()
	if err := token.Set(jwt.ExpirationKey, exp.Add(time.Minute*15)); nil != err {

		return "", fmt.Errorf("unable to set auth token expiration key: %w", err)
	}

	if err := token.Set(jwt.IssuedAtKey, time.Now()); nil != err {
		return "", fmt.Errorf("unable to set auth token issued_at key: %w", err)
	}

	if err := token.Set(jwt.JwtIDKey, tokenID.String()); nil != err {
		return "", fmt.Errorf("unable to set auth token jwt_id key: %w", err)
	}

	nbf := time.Now()
	if err := token.Set(jwt.NotBeforeKey, nbf); nil != err {
		return "", fmt.Errorf("unable to set auth token not_before key: %w", err)
	}

	if err := token.Set(jwt.SubjectKey, username); nil != err {
		return "", fmt.Errorf("unable to set auth token subject key: %w", err)
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
