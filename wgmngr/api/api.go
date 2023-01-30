package api

import (
	"database/sql"

	"github.com/xeptore/wireguard-manager/wgmngr/wg"
)

type Handler struct {
	tokenSecret  []byte
	db           *sql.DB
	wgServerConf *wg.WGServerConfig
}

func NewHandler(tokenSecret []byte, db *sql.DB, wgServerConf *wg.WGServerConfig) Handler {
	return Handler{tokenSecret, db, wgServerConf}
}
