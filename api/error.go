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

func ErrInvalidID() Error {
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

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid JSON request",
	}
}


func ErrNotResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " resource not found",
	}
}
