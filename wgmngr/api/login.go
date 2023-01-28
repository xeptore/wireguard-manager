package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/table"
	"github.com/xeptore/wireguard-manager/wgmngr/password"
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
	if err := token.Set(jwt.ExpirationKey, exp.Add(time.Hour*24*7)); nil != err {

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

var (
	ErrInternal     error = errors.New("internal error")
	ErrUserNotFound error = errors.New("user not found")
	ErrInvalidCreds error = errors.New("username or password is invalid")
)

func getUsernamePassword(ctx context.Context, db *sql.DB, username string) ([]byte, error) {
	var u m.Users
	err := t.Users.SELECT(t.Users.Password).WHERE(t.Users.Name.EQ(mysql.String(username))).LIMIT(1).QueryContext(ctx, db, &u)
	if nil != err {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u.Password, nil
}

func (s *Handler) Login(ctx context.Context, username, passwd string) (string, error) {
	storedPasswd, err := getUsernamePassword(ctx, s.db, username)
	if nil != err {
		log.Error().Err(err).Msg("failed to retrieve user stored password")
		return "", ErrInternal
	}

	matches, err := password.Compare(storedPasswd, []byte(passwd))
	if nil != err {
		log.Error().Err(err).Msg("failed to compare entered password and stored password")
		return "", ErrInternal
	}
	if !matches {
		return "", ErrInvalidCreds
	}

	token, err := generateToken(s.tokenSecret, username)
	if nil != err {
		log.Error().Err(err).Msg("failed to generate token")
		return "", ErrInternal
	}
	return token, nil
}
