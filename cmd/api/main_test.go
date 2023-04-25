package main

import (
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/jsonlog"
	"github.com/Skaifai/gophers-online-store/internal/mailer"
	_ "github.com/lib/pq"
	"os"
)

var testApp = newTestApplication()

func newTestApplication() *application {
	var cfg = config{
		port: 7000,
		env:  "test",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}{
			dsn:          "postgres://postgres:0000@localhost/gophers?sslmode=disable",
			maxOpenConns: 25, maxIdleConns: 25, maxIdleTime: "25m",
		},
		limiter: struct {
			enabled bool
			rps     float64
			burst   int
		}{
			enabled: true, rps: 2, burst: 4,
		},
		smtp: struct {
			host     string
			port     int
			username string
			password string
			sender   string
		}{
			host: "smtp.office365.com", port: 587,
			username: "211121@astanait.edu.kz",
			password: "dimok456",
			sender:   "Gophers Team 211121@astanait.edu.kz",
		},
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	//defer db.Close()
	//logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	//err = app.serve()
	//if err != nil {
	//	logger.PrintFatal(err, nil)
	//}
	return app
}
