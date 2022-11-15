package repository

import (
	"avito-job/src/internal/domain"
	"context"
)

type Repository interface {
	GetBalance(ctx context.Context, userId uint) (domain.Money, error)
	ReplenishBalance(ctx context.Context, userId uint, amount domain.Money) error
	ReserveMoney(ctx context.Context, userId uint, amount domain.Money, serviceId uint, orderId uint, description string) error
	RecognizeRevenue(ctx context.Context, userId uint, amount domain.Money, serviceId uint, orderId uint) error
}
