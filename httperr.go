package httperr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Status interface {
	error
	StatusCode() int
	Write(w http.ResponseWriter)
}

type StatusError struct {
	Code int
	Err  error
}

func Handle(w http.ResponseWriter, err Status) bool {
	if err == nil {
		return false
	}
	err.Write(w)
	return true
}

func HandleGeneric(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	return Handle(w, Server(err))
}

func (s StatusError) Error() string {
	return s.Err.Error()
}

func (s StatusError) StatusCode() int {
	return s.Code
}

func (s StatusError) Write(w http.ResponseWriter) {
	b, _ := json.Marshal(map[string]string{"message": s.Error()})
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(s.Code)
	w.Write(b)
}

func new(code int, err error) StatusError {
	return StatusError{
		Code: code,
		Err:  err,
	}
}

func NewFromString(code int, err string) StatusError {
	return new(code, fmt.Errorf(err))
}

func Server(err error) StatusError {
	return new(http.StatusInternalServerError, err)
}

func NotFound(msg string) StatusError {
	return new(http.StatusNotFound, fmt.Errorf(msg))
}

func Forbidden(msg string) StatusError {
	return new(http.StatusForbidden, fmt.Errorf(msg))
}

func Unauthorized(msg string) StatusError {
	return new(http.StatusUnauthorized, fmt.Errorf(msg))
}

func BadRequest(err error) StatusError {
	return new(http.StatusBadRequest, err)
}
