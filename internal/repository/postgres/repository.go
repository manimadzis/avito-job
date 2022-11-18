package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/manimadzis/avito-job/internal/domain"
	"github.com/manimadzis/avito-job/internal/repository"
	"github.com/manimadzis/avito-job/pkg/logging"
)

type repo struct {
	db     *sqlx.DB
	logger logging.Logger
}

func (r repo) GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error) {
	r.logger.Tracef("GetBalance(%v, %#v)", ctx, *dto)
	var amount domain.Money
	row := r.db.QueryRowContext(ctx, "select get_balance($1)", dto.UserId)
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
		r.logger.Errorf("GetBalance error: %v", err)
	}
	return amount, err
}

func (r repo) ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error {
	r.logger.Tracef("ReplenishBalance(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL replenish_balance($1, $2, $3)", dto.UserId, dto.Amount.String(), dto.Description)
	if err != nil {
		r.logger.Errorf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error {
	r.logger.Tracef("ReserveMoney(%v, %#v)", ctx, *dto)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("ReserveMoney error: %v", err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "CALL reserve_money($1, $2, $3, $4, $5)",
		dto.UserId,
		dto.Amount.String(),
		dto.ServiceId,
		dto.OrderId,
		dto.Description)
	if err != nil {
		r.logger.Errorf("Reserve money error: %v", err)
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "no_data_found" {
				return repository.ErrUnknownUser
			} else if pqerr.Message == "NOT_ENOUGH_MONEY" {
				return repository.ErrNotEnoughMoney
			} else if pqerr.Code.Name() == "unique_violation" {
				return repository.ErrTransactionAlreadyExists
			}
		}
	}

	if dto.ServiceName != "" {
		_, err := tx.ExecContext(ctx, "CALL add_service($1, $2)",
			dto.ServiceId,
			dto.ServiceName)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("ReserveMoney commit failed: %v", err)
		return err
	}

	return nil
}

func (r repo) RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error {
	r.logger.Tracef("RecognizeRevenue(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL recognize_revenue($1, $2, $3, $4)",
		dto.UserId,
		dto.Amount.String(),
		dto.ServiceId,
		dto.OrderId)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Message == "UNKNOWN_TRANSACTION" {
				return repository.ErrUnknownTransaction
			}
		}
		r.logger.Errorf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) CancelTransaction(ctx context.Context, dto *domain.CancelTransactionDTO) error {
	r.logger.Tracef("CancelTransaction(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL cancel_transaction($1, $2, $3, $4)",
		dto.UserId,
		dto.Amount.String(),
		dto.ServiceId,
		dto.OrderId)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Message == "UNKNOWN_TRANSACTION" {
				return repository.ErrUnknownTransaction
			}
		}
		r.logger.Errorf("CancelTransaction error: %v", err)
	}
	return err
}

func (r repo) GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error) {
	r.logger.Tracef("GetMonthlyReportPath(%v, %#v)", ctx, *dto)
	rows, err := r.db.QueryxContext(ctx, "SELECT * FROM get_month_report($1, $2)",
		dto.Month,
		dto.Year)
	if err != nil {
		r.logger.Errorf("GetMonthlyReportPath error: %v", err)
		return nil, err
	}

	var row domain.MonthlyReportRow
	var report domain.MonthlyReport
	for rows.Next() {
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
		report = append(report, row)
	}
	return report, nil
}

func (r repo) GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error) {
	r.logger.Tracef("GetMonthlyReportPath(%v, %#v)", ctx, *dto)
	var rows *sqlx.Rows
	var err error
	if dto.SortBy == "" || dto.SortBy == domain.GetHistoryDTOSortByTimestamp {
		rows, err = r.db.QueryxContext(ctx, "SELECT * FROM get_history_sorted_by_timestamp($1, $2, $3, $4)",
			dto.UserId,
			dto.Offset,
			dto.Limit,
			dto.Reverse)
	} else if dto.SortBy == domain.GetHistoryDTOSortByAmount {
		rows, err = r.db.QueryxContext(ctx, "SELECT * FROM get_history_sorted_by_amount($1, $2, $3, $4)",
			dto.UserId,
			dto.Offset,
			dto.Limit,
			dto.Reverse)
	}
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "no_data_found" && pqerr.Message == "UNKNOWN_USER" {
				return nil, repository.ErrUnknownUser
			}
		}
		r.logger.Errorf("GetMonthlyReportPath error: %v", err)
		return nil, err
	}
	var row domain.HistoryRow
	var history domain.History
	for rows.Next() {
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
		history = append(history, row)
	}
	return history, nil
}

func NewRepository(db *sqlx.DB, logger logging.Logger) repository.Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}
