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
	fmt.Println(cfg)

	// DB connection
	db, err := sql.Open("sqlite3", path.Join(cfg.Storage.Path, cfg.Storage.Filename))
	if err != nil {
		fmt.Println(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Persistent Storage
	persistenrStorage, err := sqlite.New(db)
	if err != nil {
		fmt.Println(err)
	}

	// Inmemory Storage
	inmemoryStorage := inmemory.New()

	// Service
	service, err := service.New(persistenrStorage, inmemoryStorage)
	if err != nil {
		fmt.Println(err)
	}

	// Presenter
	presenter := presenter.New(service)

	// Terminal user interface.
	t := tui.New(service, presenter)
	if err := t.Run(); err != nil {
		t.Stop()
	}
}
