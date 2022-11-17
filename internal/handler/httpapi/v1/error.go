package v1

import "fmt"

var (
	ErrInternalServer = fmt.Errorf("internal server error")
	ErrInvalidUserId  = fmt.Errorf("invalid user_id")
	ErrYear           = fmt.Errorf("invalid year")
	ErrMonth          = fmt.Errorf("invalid year")
)

type ErrorResponse struct {
	Msg string `json:"msg"`
}
