package app

import (
	"avito-job/internal/config"
	"avito-job/internal/repository"
	"avito-job/internal/repository/postgres"
	"avito-job/internal/server"
	"avito-job/internal/service"
	dbclient "avito-job/pkg/dbclient/postgres"
	"avito-job/pkg/logging"
	"context"
	"github.com/jmoiron/sqlx"
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
	a.service = service.NewService(a.repo, a.logger)
	a.server = server.NewServer(&server.Config{
		Host: a.config.ServerHost,
		Port: a.config.ServerPort,
	}, a.service, a.logger)

	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
