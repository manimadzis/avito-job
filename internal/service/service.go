package service

import (
	"avito-job/internal/repository"
	"avito-job/pkg/logging"
)

type Service interface {
}

type service struct {
	repo   repository.Repository
	logger logging.Logger
}

func NewService(repo repository.Repository, logger logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}
