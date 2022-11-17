package postgres

import (
	"avito-job/internal/domain"
	repository2 "avito-job/internal/repository"
	"avito-job/pkg/logging"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type repo struct {
	db     *sqlx.DB
	logger logging.Logger
}

func (r repo) GetBalance(ctx context.Context, dto *domain.GetBalanceDTO) (domain.Money, error) {
	r.logger.Tracef("GetBalance(%v, %#v)", ctx, *dto)
	var amount float64
	row := r.db.QueryRowContext(ctx, "select get_balance($1)", dto.UserId)
	err := row.Err()
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "no_data_found" {
				r.logger.Debugf("GetBalance error: %v", repository2.ErrUnknownUser)
				return domain.Money(0), repository2.ErrUnknownUser
			}
		}
		r.logger.Debugf("GetBalance error: %v", repository2.ErrUnknownUser)
		return domain.Money(0), err
	}
	err = row.Scan(&amount)
	if err != nil {
		r.logger.Errorf("GetBalance error: %v", repository2.ErrUnknownUser)
	}
	return domain.Float64ToMoney(amount), err
}

func (r repo) ReplenishBalance(ctx context.Context, dto *domain.ReplenishBalanceDTO) error {
	r.logger.Tracef("ReplenishBalance(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL replenish_balance($1, $2)", dto.UserId, dto.Amount.String())
	if err != nil {
		r.logger.Errorf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) ReserveMoney(ctx context.Context, dto *domain.ReserveMoneyDTO) error {
	r.logger.Tracef("ReserveMoney(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL reserve_money($1, $2, $3, $4, $5)",
		dto.UserId,
		dto.Amount.String(),
		dto.ServiceId,
		dto.OrderId,
		dto.Description)
	if err != nil {
		r.logger.Errorf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) RecognizeRevenue(ctx context.Context, dto *domain.RecognizeRevenueDTO) error {
	r.logger.Tracef("RecognizeRevenue(%v, %#v)", ctx, *dto)
	_, err := r.db.ExecContext(ctx, "CALL recognize_revenue($1, $2, $3, $4)",
		dto.UserId,
		dto.Amount.String(),
		dto.ServiceId,
		dto.OrderId)
	if err != nil {
		r.logger.Errorf("RecognizeRevenue error: %v", err)
	}
	return err
}

func (r repo) GetMonthlyReport(ctx context.Context, dto *domain.GetMonthlyReportDTO) (domain.MonthlyReport, error) {
	r.logger.Tracef("GetMonthlyReport(%v, %#v)", ctx, *dto)
	rows, err := r.db.QueryxContext(ctx, "SELECT * FROM get_month_report($1, $2)",
		dto.Month,
		dto.Year)
	if err != nil {
		r.logger.Errorf("GetMonthlyReport error: %v", err)
		return nil, err
	}
	var tmprow struct {
		ServiceName string `json:"service_name" db:"service_name"`
		Revenue     string `json:"revenue" db:"revenue"`
	}
	var row domain.MonthlyReportRow
	var report domain.MonthlyReport
	for rows.Next() {
		if err := rows.StructScan(&tmprow); err != nil {
			return nil, err
		}
		m, err := domain.StringToMoney(tmprow.Revenue)
		if err != nil {
			r.logger.Errorf("Can't convert Revernue to domain.Money")
			continue
		}
		row.ServiceName = tmprow.ServiceName
		row.Revenue = m
		report = append(report, row)
	}
	return report, nil
}

func (r repo) GetHistory(ctx context.Context, dto *domain.GetHistoryDTO) (domain.History, error) {
	r.logger.Tracef("GetMonthlyReport(%v, %#v)", ctx, *dto)
	var rows *sqlx.Rows
	var err error
	if dto.SortBy == "" || dto.SortBy == domain.GetHistoryDTOSortByTimestamp {
		rows, err = r.db.QueryxContext(ctx, "SELECT * FROM get_history_sorted_by_timestamp($1, $2, $3, $4)",
			dto.UserId,
			dto.Offset,
			dto.Limit,
			dto.Revers)
	} else if dto.SortBy == domain.GetHistoryDTOSortByAmount {
		rows, err = r.db.QueryxContext(ctx, "SELECT * FROM get_history_sorted_by_amount($1, $2, $3, $4)",
			dto.UserId,
			dto.Offset,
			dto.Limit,
			dto.Revers)
	}
	if err != nil {
		r.logger.Errorf("GetMonthlyReport error: %v", err)
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

func NewRepository(db *sqlx.DB, logger logging.Logger) repository2.Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}
