package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xeptore/wireguard-manager/wgmngr/api"
	"github.com/xeptore/wireguard-manager/wgmngr/env"
	"github.com/xeptore/wireguard-manager/wgmngr/migration"
	_ "github.com/xeptore/wireguard-manager/wgmngr/migration/seeds"
	"github.com/xeptore/wireguard-manager/wgmngr/wg"
)

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	if len(os.Args) != 2 {
		log.Fatal().Msg("server wireguard config file path argument is required")
	}
	wgServerConfigFilePath := os.Args[1]

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

	wgServerConf, err := wg.LoadConfig(ctx, wgServerConfigFilePath)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to load wireguard server configuration")
	}

	handler := api.NewHandler(
		[]byte(env.MustGet("AUTH_TOKEN_SECRET")),
		db,
		wgServerConf,
	)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowCredentials = true
	cfg.AllowOrigins = []string{"http://localhost:3000"}
	engine.Use(cors.New(cfg))

	engine.POST("auth/login", login(&handler))
	engine.GET("auth/check", isAuthenticated(&handler))
	engine.POST("peers", createPeer(&handler))
	engine.GET("peers/:id", getPeer(&handler))
	engine.GET("peers", getPeers(&handler))

	addr := ":8080"
	log.Info().Str("addr", addr).Msg("starting server...")
	if err := engine.Run(addr); nil != err {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func isMore(r io.Reader) bool {
	var buf [1]byte
	n, err := r.Read(buf[:])
	return !(errors.Is(err, io.EOF) && n == 0)
}

var (
	ErrRequestBodyTooLarge = errors.New("request body too large")
	ErrInvalidJSONBody     = errors.New("invalid json request body")
)

func parseJsonLimitedReader(limit int64, r io.ReadCloser, w http.ResponseWriter, v any) error {
	decoder := json.NewDecoder(io.LimitReader(r, limit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); nil != err {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return ErrInvalidJSONBody
	}

	if isMore(r) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return ErrRequestBodyTooLarge
	}

	return nil
}

func login(h *api.Handler) func(c *gin.Context) {
	const reqBodyLimit = len(`{"username":"","password":""}`) + 128 + 64

	type Form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	validateForm := func(f *Form) error {
		if len(f.Password) == 0 {
			return errors.New("password is required")
		}

		if len(f.Password) > 128 {
			return errors.New("password cannot be longer than 128 characters")
		}

		if len(f.Username) == 0 {
			return errors.New("username is required")
		}

		if len(f.Username) != len(strings.TrimSpace(f.Username)) {
			return errors.New("invalid username")
		}

		if len(f.Username) > 64 {
			return errors.New("username cannot be longer than 64 characters")
		}

		return nil
	}

	return func(c *gin.Context) {
		defer func() {
			if err := c.Request.Body.Close(); nil != err {
				log.Error().Err(err).Msg("failed to close request body")
			}
		}()

		var f Form
		if err := parseJsonLimitedReader(int64(reqBodyLimit), c.Request.Body, c.Writer, &f); nil != err {
			return
		}

		if err := validateForm(&f); nil != err {
			c.Writer.WriteHeader(http.StatusUnprocessableEntity)
			c.Writer.Write([]byte(err.Error()))
			return
		}

		token, err := h.Login(c.Request.Context(), f.Username, f.Password)
		if nil != err {
			if errors.Is(err, api.ErrInvalidCreds) || errors.Is(err, api.ErrUserNotFound) {
				rand.Seed(time.Now().Unix())
				time.Sleep(time.Duration(1000+rand.Int31n(200)+300) * time.Millisecond)
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.SetCookie("auth", token, int(time.Duration(time.Minute*15).Seconds()), "/", "localhost", false, true)
		c.Writer.WriteHeader(http.StatusCreated)
	}
}

func createPeer(h *api.Handler) func(c *gin.Context) {
	const reqBodyLimit = len(`{"name":"","description":""}`) + 256 + 10_000

	type Form struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	validateForm := func(f *Form) error {
		if len(f.Name) == 0 {
			return errors.New("name is required")
		}

		if len(f.Name) > 256 {
			return errors.New("name cannot be longer than 256 characters")
		}

		if len(f.Name) != len(strings.TrimSpace(f.Name)) {
			return errors.New("invalid name")
		}

		if len(f.Description) > 10000 {
			return errors.New("description cannot be longer than 10,000 characters")
		}

		if len(f.Description) != len(strings.TrimSpace(f.Description)) {
			return errors.New("invalid description")
		}

		return nil
	}

	return func(c *gin.Context) {
		defer func() {
			if err := c.Request.Body.Close(); nil != err {
				log.Error().Err(err).Msg("failed to close request body")
			}
		}()

		authCookie, err := c.Cookie("auth")
		if nil != err {
			if errors.Is(err, http.ErrNoCookie) {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		tokenClaims, err := h.ParseVerifyToken(authCookie)
		if nil != err {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		var f Form
		if err := parseJsonLimitedReader(int64(reqBodyLimit), c.Request.Body, c.Writer, &f); nil != err {
			return
		}

		if err := validateForm(&f); nil != err {
			c.Writer.WriteHeader(http.StatusUnprocessableEntity)
			c.Writer.Write([]byte(err.Error()))
			return
		}

		configID, err := h.CreatePeerConfig(c.Request.Context(), api.CreatePeerConfigReq{
			Name:        f.Name,
			Description: f.Description,
			ResellerID:  tokenClaims.UserID,
		})
		if nil != err {
			if errors.Is(err, api.ErrInvalidIPv4Address) || errors.Is(err, api.ErrInvalidIPv6Address) {
				log.Error().Err(err).Msg("invalid ip address")
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			if errors.Is(err, api.ErrNoMoreIPv4Address) || errors.Is(err, api.ErrNoMoreIPv6Address) {
				log.Error().Err(err).Msg("no more ip address to associate")
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Error().Err(err).Msg("failed to create peer config by unknown error")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusCreated)
		c.Writer.Write([]byte(configID))
	}
}

func getPeer(h *api.Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := c.Request.Body.Close(); nil != err {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("failed to close request body")
		}

		authCookie, err := c.Cookie("auth")
		if nil != err {
			if errors.Is(err, http.ErrNoCookie) {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		tokenClaims, err := h.ParseVerifyToken(authCookie)
		if nil != err {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, ok := c.Params.Get("id")
		if !ok {
			c.Writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		configContent, err := h.GetPeerConfig(c.Request.Context(), api.GetPeerConfigReq{
			ResellerID: tokenClaims.UserID,
			ConfigID:   id,
		})
		if nil != err {
			log.Error().Err(err).Msg("failed to get peer config")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(configContent)
	}
}

func getPeers(h *api.Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := c.Request.Body.Close(); nil != err {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("failed to close request body")
		}

		configsContent, err := h.GetActivePeerConfigs(c.Request.Context())
		if nil != err {
			log.Error().Err(err).Msg("failed to get peer configs")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(configsContent)
	}
}

func isAuthenticated(h *api.Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := c.Request.Body.Close(); nil != err {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("failed to close request body")
		}

		authCookie, err := c.Cookie("auth")
		if nil != err {
			if errors.Is(err, http.ErrNoCookie) {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ok, err := h.IsAuthenticated(c.Request.Context(), authCookie)
		if nil != err {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		c.Writer.WriteHeader(http.StatusOK)
	}
}
