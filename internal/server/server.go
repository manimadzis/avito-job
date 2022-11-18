package server

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/manimadzis/avito-job/internal/handler/httpapi/v1"
	"github.com/manimadzis/avito-job/internal/service"
	"github.com/manimadzis/avito-job/pkg/logging"
	"net/http"
	"time"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type server struct {
	logger     logging.Logger
	config     *Config
	service    service.Service
	httpServer *http.Server
}

func NewServer(config *Config, service service.Service, logger logging.Logger) Server {
	return &server{
		logger:  logger,
		service: service,
		config:  config,
		httpServer: &http.Server{
			Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
			Handler: v1.NewHandler(&v1.Config{
				Directory: config.FileServerDirectory,
				ServerURI: fmt.Sprintf("%s:%s", config.Host, config.Port),
			}, httprouter.New(), service, logger),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}
}

func (s *server) ListenAndServe() error {
	s.logger.Infof("Starting server on %s:%s", s.config.Host, s.config.Port)
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
