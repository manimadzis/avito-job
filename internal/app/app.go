package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/manimadzis/avito-job/internal/config"
	"github.com/manimadzis/avito-job/internal/repository"
	"github.com/manimadzis/avito-job/internal/repository/postgres"
	"github.com/manimadzis/avito-job/internal/server"
	"github.com/manimadzis/avito-job/internal/service"
	dbclient "github.com/manimadzis/avito-job/pkg/dbclient/postgres"
	"github.com/manimadzis/avito-job/pkg/logging"
)

type App struct {
	config  *config.Config
	logger  logging.Logger
	db      *sqlx.DB
	repo    repository.Repository
	service service.Service
	server  server.Server
}

func NewApp(config *config.Config, logger logging.Logger) *App {
	return &App{
		config: config,
		logger: logger,
	}
}

func (a *App) Start() error {
	a.logger.Info("Starting app...")
	var err error
	a.db, err = dbclient.New(dbclient.Config{
		Host:     a.config.DBHost,
		Port:     a.config.DBPort,
		Username: a.config.DBUsername,
		Password: a.config.DBPassword,
		Database: a.config.DatabaseName,
	})
	if err != nil {
		return err
	}
	defer a.db.Close()
	a.repo = postgres.NewRepository(a.db, a.logger)
	a.service = service.NewService(&service.Config{FileServerDirectory: a.config.FileServerDirectory}, a.repo, a.logger)
	a.server = server.NewServer(&server.Config{
		Host: a.config.ServerHost,
		Port: a.config.ServerPort,
	}, a.service, a.logger)

	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
