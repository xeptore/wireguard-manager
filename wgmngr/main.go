package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	tz, isTzSet := os.LookupEnv("TZ")
	if !isTzSet || tz != "UTC" {
		log.Fatal().Msg("TZ environment variable must set to UTC")
	}

	// https://github.com/go-sql-driver/mysql/#parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=false&loc=UTC&parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_DATABASE"),
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

	handler := api.NewHandler(nil, db)
	router := httprouter.New()
	router.POST("/auth/login", login(&handler))

	addr := ":8080"
	log.Info().Str("addr", addr).Msg("starting server")
	if err := http.ListenAndServe(addr, router); nil != err {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func login(h *api.Handler) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token, err := h.Login(r.Context(), "", "")
		if errors.Is(err, api.ErrInvalidCreds) || errors.Is(err, api.ErrUserNotFound) {
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
