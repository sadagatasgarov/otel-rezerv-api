package api

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// Error implemen edir Error interfeysini
func (e Error) Error() string {
	return e.Err
}


func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Err:  msg,
	}
}

func ErrInvalid() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}
