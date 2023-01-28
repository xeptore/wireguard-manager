package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xeptore/wireguard-manager/wgmngr/api"
	"github.com/xeptore/wireguard-manager/wgmngr/env"
	"github.com/xeptore/wireguard-manager/wgmngr/migration"
	_ "github.com/xeptore/wireguard-manager/wgmngr/migration/seeds"
)

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	if err := godotenv.Load(); nil != err {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic().Err(err).Msg("unexpected error while loading .env file")
		}
		log.Warn().Msg(".env file not found")
	}

	tz := env.MustGet("TZ")
	if tz != "UTC" {
		log.Fatal().Msg("TZ environment variable must be set to UTC")
	}

	// https://github.com/go-sql-driver/mysql/#parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=false&loc=UTC&parseTime=true",
		env.MustGet("DB_USERNAME"),
		env.MustGet("DB_PASSWORD"),
		env.MustGet("DB_ADDRESS"),
		env.MustGet("DB_DATABASE"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database connection")
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err := db.PingContext(ctx); nil != err {
		log.Fatal().Err(err).Msg("failed to ping the database")
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetTableName("migrations")
	goose.SetBaseFS(migration.FS)

	if err := goose.SetDialect("mysql"); nil != err {
		log.Fatal().Err(err).Msg("failed to set goose sql dialect to mysql")
	}

	log.Trace().Msg("executing database migrations...")
	if err := goose.Up(db, "scripts"); nil != err {
		log.Fatal().Err(err).Msg("failed to run migration scripts using goose")
	}
	log.Info().Msg("database migrations executed")

	handler := api.NewHandler([]byte(env.MustGet("AUTH_TOKEN_SECRET")), db)
	router := httprouter.New()
	router.POST("/auth/login", login(&handler))

	addr := ":8080"
	log.Info().Str("addr", addr).Msg("starting server")
	if err := http.ListenAndServe(addr, router); nil != err {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func isMore(r io.Reader) bool {
	var buf [1]byte
	n, err := r.Read(buf[:])
	return !(errors.Is(err, io.EOF) && n == 0)
}

type ErrorTodo struct{}

func (ErrorTodo) Error() string {
	return ""
}

var ErrTODO = ErrorTodo{}

func parseJsonLimitedReader(r io.ReadCloser, w http.ResponseWriter, limit int64, v any) error {
	decoder := json.NewDecoder(io.LimitReader(r, limit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); nil != err {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return ErrTODO
	}

	if isMore(r) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return ErrTODO
	}

	return nil
}

func login(h *api.Handler) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	const reqBodyLimit = len(`{"username":"","password":""}`) + 128 + 64

	type Form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := r.Body.Close(); nil != err {
				log.Error().Err(err).Msg("failed to close request body")
			}
		}()

		var f Form
		if err := parseJsonLimitedReader(r.Body, w, int64(reqBodyLimit), &f); errors.Is(err, ErrTODO) {
			return
		}

		token, err := h.Login(r.Context(), f.Username, f.Password)
		if errors.Is(err, api.ErrInvalidCreds) || errors.Is(err, api.ErrUserNotFound) {
			time.Sleep(1500 * time.Millisecond)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if errors.Is(err, api.ErrInternal) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(token))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(token))
	}
}
