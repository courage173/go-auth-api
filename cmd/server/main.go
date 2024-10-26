package main

import (
	"context"
	"strconv"

	"database/sql"

	"flag"

	dbx "github.com/go-ozzo/ozzo-dbx"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"

	"github.com/go-ozzo/ozzo-routing/v2/cors"

	_ "github.com/lib/pq"

	"github.com/courage173/quiz-api/pkg/accesslog"

	"github.com/courage173/quiz-api/pkg/dbcontext"

	"github.com/courage173/quiz-api/internal/errors"
	"github.com/courage173/quiz-api/internal/healthcheck"

	"github.com/courage173/quiz-api/internal/auth"

	"github.com/courage173/quiz-api/internal/users"

	"github.com/courage173/quiz-api/pkg/log"

	"github.com/courage173/quiz-api/pkg/utils"

	"net/http"

	"os"

	"time"

	"github.com/joho/godotenv"
)

// Version indicates the current version of the application.
var Version = "1.0.0"


func main(){
	flag.Parse()
	godotenv.Load()

  // create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	 err := godotenv.Load()
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	databaseUrl := os.Getenv("DB_URL")

	if databaseUrl == "" {
        logger.Errorf("DB_URL is required")
		os.Exit(-1)
    }
	// connect to the database
	db, err := dbx.MustOpen("postgres", databaseUrl)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	// build HTTP server
	portStr := os.Getenv("HTTP_PORT")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        logger.Errorf("Invalid port number: %s", portStr)
		os.Exit(-1)
    }

    // Ensure the port is within the valid range
    if port < 1 || port > 65535 {
        logger.Errorf("Port number out of range: %d", port)
		os.Exit(-1)
    }

	address := "localhost:" + strconv.Itoa(port)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db)),
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
	     errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	 healthcheck.RegisterHandlers(router, Version)

	 rg := router.Group("/v1")

	// authHandler := auth.Handler(cfg.JWTSigningKey)

	// album.RegisterHandlers(rg.Group(""),
	// 	album.NewService(album.NewRepository(db, logger), logger),
	// 	authHandler, logger,
	// )
	jwtSigningKey := utils.GetEnv("JWT_SIGNING_KEY")
	jwtExpiration := utils.GetEnv("JWT_EXPIRY_")

	if jwtExpiration == "" {
		jwtExpiration = "3600"
	}

	expiry, _ := strconv.Atoi(jwtExpiration);

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(jwtSigningKey, expiry , logger, users.NewRepository(db, logger)),
		logger,
	)

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}