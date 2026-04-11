package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/re0udjat/magic-shop/internal/data"
	"github.com/re0udjat/magic-shop/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxConns     int
		minIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type app struct {
	config config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	// Flags for app
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Envrionment (dev|stag|prod)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/magical_shop?sslmode=disable", "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxConns, "db-max-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minIdleConns, "db-min-idle-conns", 5, "PostgreSQL min idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter max burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "242dfbdff3aa7a", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "daa4a46ad0924e", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Magical Shop <no-reply@magical-shop.cuongnlt.com>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	// Structured logger which writes log entries to stdout stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open the connection pool
	dbpool, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to dbpool.Close() so that the connection pool is closed before the main()
	// function exits
	defer dbpool.Close()
	logger.Info("database connection pool established")

	// Publish a new "version" variable in the expvar handler containing our app version
	// number (currently the constant "1.0.0")
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// Publish the db connection pool statistics
	expvar.Publish("database", expvar.Func(func() any {
		return dbpool.Stat()
	}))

	// Publish the current Unix timestamp
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &app{
		config: cfg,
		logger: logger,
		models: data.NewModels(dbpool),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	/* Declare a server:
	*	- IdleTimeout: max duration for keep-alive connection
	*	- ReadTimeout: max duration for reading the request headers and body => Mitigate the risk from slow-client attacks
	*	- WriteTimeout: max duration for writing the response body
	 */
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	// Create a connection pool
	dbpool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Max number of opens connections in the pool (in-use + idle connections)
	dbpool.Config().MaxConns = int32(cfg.db.maxConns)

	// Min number of idle connections in the pool
	dbpool.Config().MinIdleConns = int32(cfg.db.minIdleConns)

	// Max connection idle time
	dbpool.Config().MaxConnIdleTime = cfg.db.maxIdleTime

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to verify the connection
	err = dbpool.Ping(ctx)
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
