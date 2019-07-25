package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

type Reply interface {
	StatusCode() int
	Encode(io.Writer) error
}

type reply struct {
	value      interface{}
	statusCode int
}

func NewReply(value interface{}, statusCode int) Reply {
	return &reply{
		value:      value,
		statusCode: statusCode,
	}
}

func NewOKReply(value interface{}) Reply {
	return NewReply(value, http.StatusOK)
}

func (r *reply) StatusCode() int {
	return r.statusCode
}

func (r *reply) Encode(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r.value)
}
