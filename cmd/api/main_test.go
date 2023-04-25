package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/jsonlog"
	"github.com/Skaifai/gophers-online-store/internal/mailer"
	"os"
	"strconv"
)

//var testingApplication *application = SetupApplication(cfg, jsonlog.New(os.Stdout, jsonlog.LevelInfo), openDB(cfg))

func SetupApplication(cfg config, logger *jsonlog.Logger, db *sql.DB) *application {
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	return app
}

func SetupConfig() config {
	var cfg config
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Empty")
		port = "7000"
	}
	port_int, _ := strconv.Atoi(port)
	flag.IntVar(&cfg.port, "port", port_int, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// GetById the database connection string, aka data source name (DSN)
	flag.StringVar(&cfg.db.dsn, "db-dsn", getEnvVar("DB_URL"), "PostgreSQL DSN")

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
	flag.StringVar(&cfg.smtp.username, "smtp-username", getEnvVar("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", getEnvVar("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Gopher Team <>", "SMTP sender")

	//flag.Parse()

	return cfg
}

var testingApplication = func() *application {
	var cfg = SetupConfig()
	db, _ := openDB(cfg)
	return SetupApplication(cfg, jsonlog.New(os.Stdout, jsonlog.LevelInfo), db)
}()
