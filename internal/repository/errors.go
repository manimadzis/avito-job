package repository

import "fmt"

var (
	ErrUnknownUser              = fmt.Errorf("unknown user")
	ErrNotEnoughMoney           = fmt.Errorf("not enough money")
	ErrUnknownTransaction       = fmt.Errorf("unknown transaction")
	ErrTransactionAlreadyExists = fmt.Errorf("transaction already exists")
)
