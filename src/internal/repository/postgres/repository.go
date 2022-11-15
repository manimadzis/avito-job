package postgres

import (
	"avito-job/src/internal/domain"
	"avito-job/src/internal/repository"
	"avito-job/src/pkg/logging"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type repo struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func (r repo) GetBalance(ctx context.Context, userId uint) (domain.Money, error) {
	r.logger.Tracef("GetBalance(%v, %d)", ctx, userId)
	var amount float64
	row := r.db.QueryRowContext(ctx, "select get_balance($1)", userId)
	err := row.Err()
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "no_data_found" {
				r.logger.Debugf("GetBalance error: %v", repository.ErrUnknownUser)
				return domain.Money(0), repository.ErrUnknownUser
			}
		}
		r.logger.Debugf("GetBalance error: %v", repository.ErrUnknownUser)
		return domain.Money(0), err
	}
	err = row.Scan(&amount)
	if err != nil {
		r.logger.Debugf("GetBalance error: %v", repository.ErrUnknownUser)
	}
	return domain.Float64ToMoney(amount), err
}

func (r repo) ReplenishBalance(ctx context.Context, userId uint, amount domain.Money) error {
	r.logger.Tracef("ReplenishBalance(%v, %d, %v)", ctx, userId, amount)
	_, err := r.db.ExecContext(ctx, "CALL replenish_balance($1, $2)", userId, amount.String())
	if err != nil {
		r.logger.Debugf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) ReserveMoney(
	ctx context.Context,
	userId uint,
	amount domain.Money,
	serviceId uint,
	orderId uint,
	description string) error {
	r.logger.Tracef("ReserveMoney(%v, %v, %v, %v, %v %v)", ctx, userId, amount, serviceId, orderId, description)
	_, err := r.db.ExecContext(ctx, "CALL reserve_money($1, $2, $3, $4, $5)",
		userId,
		amount.String(),
		serviceId,
		orderId,
		description)
	if err != nil {
		r.logger.Debugf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) RecognizeRevenue(
	ctx context.Context,
	userId uint,
	amount domain.Money,
	serviceId uint,
	orderId uint) error {
	r.logger.Tracef("RecognizeRevenue(%v, %v, %v, %v, %v)", ctx, userId, amount, serviceId, orderId)
	_, err := r.db.ExecContext(ctx, "CALL recognize_revenue($1, $2, $3, $4)",
		userId,
		amount.String(),
		serviceId,
		orderId)
	if err != nil {
		r.logger.Debugf("RecognizeRevenue error: %v", err)
	}
	return err
}

func NewRepository(db *sqlx.DB, logger *logging.Logger) repository.Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}
