package api

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

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
