package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Request interface {
	Method() string
	URL() *url.URL
	Decode(v interface{}) error
	// TODO: add more methods
}

type request struct {
	httpRequest *http.Request
}

func (r *request) Method() string {
	return r.httpRequest.Method
}

func (r *request) URL() *url.URL {
	return r.httpRequest.URL
}

func (r *request) Decode(v interface{}) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
		}
	}()

	defer func() {
		errClose := r.httpRequest.Body.Close()
		if err != nil {
			return
		}

		err = fmt.Errorf("failed to close HTTP request body: %w", errClose)
	}()

	decoder := json.NewDecoder(r.httpRequest.Body)
	err = decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode JSON body: %w", err)
	}

	return nil
}
