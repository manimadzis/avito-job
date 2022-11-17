package v1

import "fmt"

var (
	ErrInvalidUserId = fmt.Errorf("invalid user_id")
	ErrYear          = fmt.Errorf("invalid year")
	ErrMonth         = fmt.Errorf("invalid year")
	ErrUnknownUser   = fmt.Errorf("unknown user")
	ErrEmptyBody     = fmt.Errorf("empty body")
	ErrEmptyJSON     = fmt.Errorf("empty body")
)

type ErrorResponse struct {
	Msg string `json:"msg"`
}
