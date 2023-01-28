package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/v1/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/v1/wgmngr/table"
	"github.com/xeptore/wireguard-manager/wgmngr/password"
)

func init() {
	goose.AddMigration(upSeedResellers, downSeedResellers)
}

var initialResellers []m.Users = []m.Users{
	{
		Name:     "Agha Mohsen",
		Username: "mohsen",
	},
}

func upSeedResellers(tx *sql.Tx) error {
	var godUser m.Users
	err := t.Users.SELECT(t.Users.ID).WHERE(t.Users.Username.EQ(mysql.String("god"))).LIMIT(1).Query(tx, &godUser)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return errors.New("could not find god user in the database")
		}
		return errors.New("failed to query for god user")
	}

	id, err := gonanoid.New(64)
	if nil != err {
		return err
	}

	password, err := password.Hash([]byte("password"))
	if nil != err {
		return err
	}

	for i := 0; i < len(initialResellers); i++ {
		v := &initialResellers[i]
		v = &m.Users{
			ID:        id,
			Name:      v.Name,
			Username:  v.Username,
			Password:  password,
			Role:      m.UsersRole_Reseller,
			CreatorID: godUser.ID,
			CreatedAt: time.Now(),
		}
		log.Debug().Dict("reseller", zerolog.Dict().Time("created-at", v.CreatedAt).Str("username", v.Username).Str("name", v.Name).Str("id", v.ID)).Msg("seeding reseller user")
	}

	res, err := t.Users.INSERT(t.Users.AllColumns).MODELS(initialResellers).Exec(tx)
	if nil != err {
		return err
	}

	rows, err := res.RowsAffected()
	if nil != err {
		return err
	}
	if l := len(initialResellers); l != int(rows) {
		return fmt.Errorf("could not insert all initial resellers. expected %d rows to be inserted got %d", l, rows)
	}

	return nil
}

func downSeedResellers(tx *sql.Tx) error {
	usernames := make([]mysql.Expression, len(initialResellers))
	for i, v := range initialResellers {
		usernames[i] = mysql.String(v.Username)
	}

	stmt := t.Users.DELETE().WHERE(t.Users.Username.IN(usernames...))
	fmt.Println(stmt.DebugSql())
	res, err := stmt.Exec(tx)
	if nil != err {
		return err
	}

	rows, err := res.RowsAffected()
	if nil != err {
		return err
	}
	if l := len(initialResellers); l != int(rows) {
		return fmt.Errorf("could not delete all initial resellers. expected %d rows to be deleted got %d", l, rows)
	}

	return nil
}
