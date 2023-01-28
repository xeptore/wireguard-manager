package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/rs/zerolog/log"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/table"
	"github.com/xeptore/wireguard-manager/wgmngr/password"
)

var (
	ErrInternal     error = errors.New("internal error")
	ErrUserNotFound error = errors.New("user not found")
	ErrInvalidCreds error = errors.New("username or password is invalid")
)

func getUsernamePassword(ctx context.Context, db *sql.DB, username string) ([]byte, error) {
	var u m.Users
	err := t.Users.SELECT(t.Users.Password).WHERE(t.Users.Username.EQ(mysql.String(username))).LIMIT(1).QueryContext(ctx, db, &u)
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
		if errors.Is(err, ErrUserNotFound) {
			return "", err
		}

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
