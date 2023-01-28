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
	ErrUserNotFound error = errors.New("user not found")
	ErrInvalidCreds error = errors.New("username or password is invalid")
)

type userLoginCreds struct {
	id       string
	role     string
	password []byte
}

func getUserLoginInfo(ctx context.Context, db *sql.DB, username string) (*userLoginCreds, error) {
	var u m.Users
	err := t.Users.SELECT(t.Users.ID, t.Users.Role, t.Users.Password).WHERE(t.Users.Username.EQ(mysql.String(username))).LIMIT(1).QueryContext(ctx, db, &u)
	if nil != err {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &userLoginCreds{u.ID, u.Role.String(), u.Password}, nil
}

func (s *Handler) Login(ctx context.Context, username, passwd string) (string, error) {
	u, err := getUserLoginInfo(ctx, s.db, username)
	if nil != err {
		if errors.Is(err, ErrUserNotFound) {
			return "", err
		}

		log.Error().Err(err).Msg("failed to retrieve user stored password")
		return "", err
	}

	matches, err := password.Compare(u.password, []byte(passwd))
	if nil != err {
		log.Error().Err(err).Msg("failed to compare entered password and stored password")
		return "", err
	}
	if !matches {
		return "", ErrInvalidCreds
	}

	token, err := generateToken(s.tokenSecret, TokenClaims{
		UserID:   u.id,
		Role:     u.role,
		Username: username,
	})
	if nil != err {
		log.Error().Err(err).Msg("failed to generate token")
		return "", err
	}
	return token, nil
}
