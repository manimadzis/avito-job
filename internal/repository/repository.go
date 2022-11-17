package repository

import (
	"avito-job/internal/domain"
	"context"
)

type Repository interface {
	// GetBalance return ErrUnknownUser if user doesn't exist
	GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error)
	ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error
	// ReserveMoney return ErrUnknownUser if user doesn't exist
	// return ErrNotEnoughMoney if user balance lower than Amount
	ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error
	//RecognizeRevenue return ErrUnknownTransaction if transaction with given fields doesn't exist
	RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error
	GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error)
	// GetHistory return ErrUnknownUser if user doesn't exist
	GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error)
}
