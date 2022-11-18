package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/manimadzis/avito-job/internal/domain"
	"github.com/manimadzis/avito-job/internal/repository"
	"github.com/manimadzis/avito-job/pkg/logging"
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
	CancelTransaction(ctx context.Context, dto *domain.CancelTransactionDTO) error
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

	for i, row := range report {
		if row.ServiceName == "" {
			report[i].ServiceName = fmt.Sprintf("Услуга №%d", row.ServiceId)
		}
	}
	if _, err := os.Stat(s.config.FileServerDirectory); os.IsNotExist(err) {
		err := os.Mkdir(s.config.FileServerDirectory, 0666)
		if err != nil {
			s.logger.Errorf("can't create %s dir", s.config.FileServerDirectory)
			return "", err
		}
	}

	filepath := filepath.Join(s.config.FileServerDirectory, fmt.Sprintf("%d-%d.csv", dto.Year, dto.Month))
	file, err := os.Create(filepath)
	if err != nil {
		s.logger.Errorf("Cant open file %s", filepath)
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
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
	if dto.Description == "" {
		dto.Description = fmt.Sprintf("Пополнение баланса")
	}
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
	if dto.ServiceName == "" {
		dto.ServiceName = fmt.Sprintf("Услуга №%d", dto.ServiceId)
	}
	if dto.Description != "" {
		dto.Description = fmt.Sprintf("Оказание услуги: %s", dto.ServiceName)
	}
	return s.repo.ReserveMoney(ctx, dto)
}

func (s *service) RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error {
	s.logger.Tracef("service.RecognizeRevenue(%v, %#v)", ctx, *dto)
	return s.repo.RecognizeRevenue(ctx, dto)
}

func (s *service) CancelTransaction(ctx context.Context, dto *domain.CancelTransactionDTO) error {
	s.logger.Tracef("service.CancelTransaction(%v, %#v)", ctx, *dto)
	return s.repo.CancelTransaction(ctx, dto)
}

func NewService(config *Config, repo repository.Repository, logger logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
		config: config,
	}
}
