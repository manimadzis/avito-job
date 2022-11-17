package service

import (
	"avito-job/internal/domain"
	"avito-job/internal/repository"
	"avito-job/pkg/logging"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type Service interface {
	GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error)
	GetMonthlyReportPath(ctx context.Context, dto *domain.GetMonthlyReportDTO) (string, error)
	ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error
	GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error)
	ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error
	RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error
}

type service struct {
	repo   repository.Repository
	logger logging.Logger
	config *Config
}

func (s *service) GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error) {
	s.logger.Tracef("service.GetBalance(%v, %#v)", ctx, *dto)
	return s.repo.GetBalance(ctx, dto)
}

func (s *service) GetMonthlyReportPath(ctx context.Context, dto *domain.GetMonthlyReportDTO) (string, error) {
	s.logger.Tracef("service.GetBalance(%v, %#v)", ctx, *dto)
	report, err := s.repo.GetMonthlyReport(ctx, dto)
	if err != nil {
		return "", err
	}
	s.logger.Debug("Report: ", report)
	filepath := filepath.Join(s.config.FileServerDirectory, fmt.Sprintf("%d-%d.csv", dto.Year, dto.Month))
	file, err := os.Create(filepath)
	if err != nil {
		s.logger.Errorf("Cant open file %s", filepath)
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	for _, row := range report {
		err = csvWriter.Write([]string{row.ServiceName, row.Revenue.String()})
		if err != nil {
			s.logger.Errorf("Can't write to file: %v", err)
			return "", nil
		}
	}
	file.Sync()

	return filepath, nil
}
func (s *service) ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error {
	s.logger.Tracef("service.ReplenishBalance(%v, %#v)", ctx, *dto)
	return s.repo.ReplenishBalance(ctx, dto)
}

func (s *service) GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error) {
	s.logger.Tracef("service.GetHistory(%v, %#v)", ctx, *dto)
	if dto.Limit == 0 {
		dto.Limit = MaxHistoryRowPerRequest
	}
	return s.repo.GetHistory(ctx, dto)
}

func (s *service) ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error {
	s.logger.Tracef("service.ReserveMoney(%v, %#v)", ctx, *dto)
	return s.repo.ReserveMoney(ctx, dto)
}

func (s *service) RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error {
	s.logger.Tracef("service.RecognizeRevenue(%v, %#v)", ctx, *dto)
	return s.repo.RecognizeRevenue(ctx, dto)
}

func NewService(config *Config, repo repository.Repository, logger logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
		config: config,
	}
}
