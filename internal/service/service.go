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
	ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error
	RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error
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
	if dto.Limit == 0 {
		dto.Limit = MaxHistoryRowPerRequest
	}
	return s.repo.GetHistory(ctx, dto)
}

func (s *service) ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error {
	return s.repo.ReserveMoney(ctx, dto)
}

func (s *service) RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error {
	return s.repo.RecognizeRevenue(ctx, dto)
}

func NewService(repo repository.Repository, logger logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}
