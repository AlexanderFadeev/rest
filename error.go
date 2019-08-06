package rest

import (
	"encoding/json"
	"github.com/AlexanderFadeev/myerrors"
	"io"
)

const errorKey = "error"

type Error interface {
	error
	Reply
}

type errorImpl struct {
	error

	statusCode int
}

func NewError(err error, statusCode int) Error {
	return &errorImpl{
		error:      err,
		statusCode: statusCode,
	}
}

func (e *errorImpl) StatusCode() int {
	return e.statusCode
}

func (e *errorImpl) Encode(w io.Writer) error {
	result := map[string]string{
		errorKey: e.Error(),
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(result)
	if err != nil {
		return myerrors.Wrap(err, "failed to encode to JSON")
	}

	return nil
}
