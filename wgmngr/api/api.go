package api

import (
	"database/sql"
)

type Handler struct {
	tokenSecret []byte
	db          *sql.DB
}

func NewHandler(tokenSecret []byte, db *sql.DB) Handler {
	return Handler{tokenSecret, db}
}
