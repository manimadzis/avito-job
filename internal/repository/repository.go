package repository

import (
	"avito-job/internal/domain"
	"context"
)

type Repository interface {
	GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error)
	ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error
	ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error
	RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error
	GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error)
	GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error)
}
