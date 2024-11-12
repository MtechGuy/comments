package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/mtechguy/comments/internal/data"
)

const appVersion = "7.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
	limiter struct {
		rps     float64 // requests per second
		burst   int     // initial requests possible
		enabled bool    // enable or disable rate limiter
	}
}

type applicationDependencies struct {
	config       serverConfig
	logger       *slog.Logger
	commentModel data.CommentModel
	userModel    data.UserModel
}

func main() {
	var setting serverConfig

	flag.IntVar(&setting.port, "port", 4000, "Server port")
	flag.StringVar(&setting.environment, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&setting.db.dsn, "db-dsn", "postgres://comments:cattalya03@localhost/comments?sslmode=disable", "PostgreSQL DSN")

	flag.Float64Var(&setting.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

	flag.IntVar(&setting.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

	flag.BoolVar(&setting.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// the call to openDB() sets up our connection pool
	db, err := openDB(setting)
	if err != nil {
		logger.Error("Database connection failed")
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("Database connection pool established")

	appInstance := &applicationDependencies{
		config:       setting,
		logger:       logger,
		userModel:    data.UserModel{DB: db},
		commentModel: data.CommentModel{DB: db},
	}

	// apiServer := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", setting.port),
	// 	Handler:      appInstance.routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	// }

	// logger.Info("Starting server", "address", apiServer.Addr, "environment", setting.environment)
	// err = apiServer.ListenAndServe()
	// logger.Error(err.Error())
	// os.Exit(1)

	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(settings serverConfig) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second)
	defer cancel()

	// let's test if the connection pool was created
	// we trying pinging it with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}
