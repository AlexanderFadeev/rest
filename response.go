package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mailru/easyjson"
)

type Response interface {
	easyjson.Marshaler
}

type PtrToResponse[T any] interface {
	Response
	*T
}

type responseWithStatus[Resp Response] struct {
	resp       Resp
	statusCode int
}

func newResponseWithStatus[Resp Response](value Resp, statusCode int) responseWithStatus[Resp] {
	return responseWithStatus[Resp]{
		resp:       value,
		statusCode: statusCode,
	}
}

func NewOKReply[Resp Response](value Resp) responseWithStatus[Resp] {
	return newResponseWithStatus(value, http.StatusOK)
}

func (r *responseWithStatus[Resp]) Encode(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r.resp)
}

//easyjson:json
type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(err error) *errorResponse {
	return &errorResponse{
		Error: err.Error(),
	}
}
