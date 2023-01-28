package migrations

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/v1/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/v1/wgmngr/table"
	"github.com/xeptore/wireguard-manager/wgmngr/password"
)

func init() {
	goose.AddMigration(upSeedGodUser, downSeedGodUser)
}

func upSeedGodUser(tx *sql.Tx) error {
	godPassword, exists := os.LookupEnv("GOD_PASSWORD")
	if !exists {
		return errors.New("GO_PASSWORD environment variable is not set")
	}
	if !password.StrongPasswordRegExp.MatchString(godPassword) {
		return errors.New("GO_PASSWORD environment variable is not strong enough")
	}

	hash, err := password.Hash(godPassword)
	if nil != err {
		return err
	}

	id, err := gonanoid.New(64)
	if nil != err {
		return err
	}

	godUser := m.Users{
		ID:        id,
		Name:      "God",
		Username:  "god",
		Password:  hash,
		CreatorID: id,
		CreatedAt: time.Now(),
		Role:      m.UsersRole_Admin,
	}
	log.Debug().Dict("god", zerolog.Dict().Str("id", godUser.ID).Time("created-at", godUser.CreatedAt)).Msg("seeding god user")

	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=0;"); nil != err {
		return err
	}
	defer func() {
		if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=1;"); nil != err {
			log.Warn().Err(err).Msg("failed to revert FOREIGN_KEY_CHECKS database option to 1")
		}
	}()

	res, err := t.Users.INSERT(t.Users.AllColumns).MODEL(godUser).Exec(tx)
	if nil != err {
		return err
	}
	rows, err := res.RowsAffected()
	if nil != err {
		return err
	}
	if rows != 1 {
		return errors.New("god user was not inserted")
	}

	return err
}

func downSeedGodUser(tx *sql.Tx) error {
	res, err := t.Users.DELETE().WHERE(t.Users.Username.EQ(mysql.String("god"))).Exec(tx)
	if nil != err {
		return err
	}
	rows, err := res.RowsAffected()
	if nil != err {
		return err
	}
	if rows != 1 {
		log.Warn().Msg("expected god user to be inserted but seems it's already deleted")
	}

	return nil
}
