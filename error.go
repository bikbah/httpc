package httpc

import (
	"errors"
	"fmt"
)

type Error struct {
	Name    string
	Err     error
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s code(%d) (%v)", e.Name, e.Message, e.Code, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Is(err error) bool {
	return errors.Is(e.Err, err)
}
