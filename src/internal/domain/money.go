package domain

import "fmt"

type Money int64

func Float64ToMoney(x float64) Money {
	return Money(100 * x)
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f", float64(m)/100)
}
