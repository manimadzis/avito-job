package domain

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Money int64

func (m *Money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var err error
	if data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("money json must be a string")
	}
	*m, err = StringToMoney(string(data[1 : len(data)-1]))
	return err
}

func (m *Money) Scan(src interface{}) error {
	var err error
	var tmpm Money
	switch src := src.(type) {
	case string:
		tmpm, err = StringToMoney(src)
	case []uint8:
		tmpm, err = StringToMoney(string(src))
	default:
		err = fmt.Errorf("src must be a string or []uint8 not %s", reflect.TypeOf(src).String())
	}
	if err != nil {
		return err
	}
	*m = tmpm
	return nil
}

func StringToMoney(s string) (Money, error) {
	var err error
	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	subs := strings.Split(s, ".")
	if len(subs) > 2 {
		return 0, fmt.Errorf("invalid string")
	}
	var intPart, fracPart int

	if len(subs) == 2 && len(subs[1]) > 2 {
		return 0, fmt.Errorf("max precision is 2")
	}

	intPart, err = strconv.Atoi(subs[0])
	if err != nil {
		return 0, fmt.Errorf("invalid string")
	}
	if len(subs) == 2 {
		fracPart, err = strconv.Atoi(subs[1])
		if err != nil {
			return 0, fmt.Errorf("invalid string")
		}
		if len(subs[1]) == 1 {
			fracPart *= 10
		}
	}

	return Money(100*intPart + fracPart), nil
}

func (m *Money) String() string {
	return fmt.Sprintf("%.2f", float64(*m)/100)
}

type MonthlyReportRow struct {
	ServiceName string `json:"service_name" db:"service_name"`
	Revenue     Money  `json:"revenue" db:"revenue"`
}
type MonthlyReport []MonthlyReportRow

type HistoryRow struct {
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Amount      Money     `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
}

type History []HistoryRow
