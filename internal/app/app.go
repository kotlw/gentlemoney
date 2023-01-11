package app

import (
	"database/sql"
	"fmt"
	"path"

	"github.com/kotlw/gentlemoney/config"
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
	"github.com/kotlw/gentlemoney/internal/tui"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	// Configuration
	cfg := config.Default()

	// Logger
	log := InitLogger(cfg.Logger.Level, cfg.Logger.Path, cfg.Logger.Filename)
	log.Debug("Config has initialized.")

	// DB connection
	p := path.Join(cfg.Storage.Path, cfg.Storage.Filename)
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		log.Fatal(fmt.Errorf("app: Run: sql.Open: %w", err))
	}
	log.WithField("path", p).Debug("sql.DB has Opened")

	if err = db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("app: Run: db.Ping: %w", err))
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(fmt.Errorf("app: Run: db.Close: %w", err))
		}
	}()

	// Persistent Storage
	persistenrStorage, err := sqlite.New(db)
	if err != nil {
		log.Fatal(fmt.Errorf("app: Run: sqlite.New: %w", err))
	}
	log.Debug("SqliteStorage has initialized.")

	// Inmemory Storage
	inmemoryStorage := inmemory.New()
	log.Debug("InmemoryStorage has initialized.")

	// Service
	service, err := service.New(persistenrStorage, inmemoryStorage)
	if err != nil {
		log.Fatal(fmt.Errorf("app: Run: service.New: %w", err))
	}
	log.Debug("Service has initialized.")

	// Presenter
	presenter := presenter.New(service)
	log.Debug("Presenter has initialized.")

	// Terminal user interface.
	t := tui.New(service, presenter)
	log.Debug("TviewApplication has initialized.")
	if err := t.Run(); err != nil {
		t.Stop()
		log.Fatal(fmt.Errorf("app: Run: t.Run: %w", err))
	}
}
