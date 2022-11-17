package service

import (
	"avito-job/internal/domain"
	"avito-job/internal/repository"
	"avito-job/pkg/logging"
	"context"
)

type Service interface {
	GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error)
	GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error)
	ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error
	GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error)
}

type service struct {
	repo   repository.Repository
	logger logging.Logger
}

func (s *service) GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error) {
	return s.repo.GetBalance(ctx, dto)
}

func (s *service) GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error) {
	return s.repo.GetMonthlyReport(ctx, dto)
}
func (s *service) ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error {
	return s.repo.ReplenishBalance(ctx, dto)
}

func (s *service) GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error) {
	return s.repo.GetHistory(ctx, dto)
}

func NewService(repo repository.Repository, logger logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}
