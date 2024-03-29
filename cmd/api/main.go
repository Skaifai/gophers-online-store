package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/jsonlog"
	"github.com/Skaifai/gophers-online-store/internal/mailer"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func getEnvVar(key string) string {
	godotenv.Load()
	return os.Getenv(key)
}

func main() {
	var cfg config
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Empty")
		port = "7000"
	}
	port_int, err := strconv.Atoi(port)
	flag.IntVar(&cfg.port, "port", port_int, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// GetById the database connection string, aka data source name (DSN)
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://aiyihnvbfdwvno:9d01118fa576ac32e3ed0f7cc7e096be04285c5f66bbfbe3c099c20e127bff7b@ec2-52-18-116-67.eu-west-1.compute.amazonaws.com:5432/dho8n1cmu6kot", "PostgreSQL DSN")

	// Set up restrictions for the database connections
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	// Set up limitations for application
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Google smtp-server connection
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.office365.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "211121@astanait.edu.kz", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "SOME_PASSWORD", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Gopher Team <211121@astanait.edu.kz>", "SMTP sender")

	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
