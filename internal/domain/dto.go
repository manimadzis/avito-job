package domain

import validation "github.com/go-ozzo/ozzo-validation"

type DTO interface {
	Validate() error
}

type GetBalanceDTO struct {
	UserId uint `json:"user_id"`
}

func (d GetBalanceDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required, validation.Min(uint(0))),
	)
}

type GetMonthlyReportDTO struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (d GetMonthlyReportDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Year, validation.Required, validation.Min(1900), validation.Max(2199)),
		validation.Field(&d.Month, validation.Required, validation.Min(1), validation.Max(12)),
	)
}

type RecognizeRevenueDTO struct {
	UserId    uint  `json:"user_id"`
	Amount    Money `json:"amount"`
	ServiceId uint  `json:"service_id"`
	OrderId   uint  `json:"order_id"`
}

func (d RecognizeRevenueDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.Amount, validation.Required, validation.Min(Money(0))),
		validation.Field(&d.ServiceId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.OrderId, validation.Required, validation.Min(uint(1))),
	)
}

type ReserveMoneyDTO struct {
	UserId      uint   `json:"user_id"`
	Amount      Money  `json:"amount"`
	ServiceId   uint   `json:"service_id"`
	OrderId     uint   `json:"order_id"`
	Description string `json:"description"`
	ServiceName string `json:"service_name"`
}

func (d ReserveMoneyDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.Amount, validation.Required, validation.Min(Money(0))),
		validation.Field(&d.ServiceId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.OrderId, validation.Required, validation.Min(uint(1))),
	)
}

type ReplenishBalanceDTO struct {
	UserId      uint   `json:"user_id"`
	Amount      Money  `json:"amount"`
	Description string `json:"description"`
}

func (d ReplenishBalanceDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.Amount, validation.Required, validation.Min(Money(0))),
	)
}

type GetHistoryDTO struct {
	UserId  uint   `json:"user_id"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	SortBy  string `json:"sort_by"`
	Reverse bool   `json:"reverse"`
}

func (d GetHistoryDTO) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required, validation.Min(uint(1))),
		validation.Field(&d.Offset, validation.Min(0)),
		validation.Field(&d.Limit, validation.Min(0)),
		validation.Field(&d.SortBy, validation.In(GetHistoryDTOSortByTimestamp, GetHistoryDTOSortByAmount)),
	)
}
