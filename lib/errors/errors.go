package errors

import (
	"demo-server/lib/util"
	"encoding/json"
)

// Error implements the error object.
type Error struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return util.Bytes2str(b)
}

// New generates a custom error.
func New(code int, detail string) error {
	return &Error{
		Code:   code,
		Detail: detail,
	}
}
