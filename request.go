package rest

import (
	"encoding/json"
	"github.com/AlexanderFadeev/myerrors"
	"net/http"
)

type Request interface {
	Decode(v interface{}) error
}

type request struct {
	httpRequest *http.Request
}

func (r *request) Decode(v interface{}) (err error) {
	defer func() {
		if err != nil {
			err = myerrors.Wrap(err, "failed to decode request")
		}
	}()

	defer func() {
		errClose := r.httpRequest.Body.Close()
		if err != nil {
			return
		}

		if errClose != nil {
			err = myerrors.Wrap(errClose, "failed to close HTTP request body")
		}
	}()

	decoder := json.NewDecoder(r.httpRequest.Body)
	err = decoder.Decode(v)
	if err != nil {
		return myerrors.Wrap(err, "failed to decode JSON body")
	}

	return nil
}
