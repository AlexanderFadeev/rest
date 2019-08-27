package rest

import (
	"encoding/json"
	"github.com/AlexanderFadeev/myerrors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

const (
	queryTag = "query"
	urlTag   = "url"
)

type Request interface {
	Decode(v interface{}) error
}

type request struct {
	httpRequest *http.Request
}

func (r *request) Decode(v interface{}) error {
	err := r.decodeBody(v)
	if err != nil {
		return myerrors.Wrap(err, "failed to decode request body")
	}

	err = r.decodeQueryString(v)
	if err != nil {
		return myerrors.Wrap(err, "failed to decode query strings")
	}

	err = r.decodeURLParams(v)
	return myerrors.Wrap(err, "failed to decode URL params")
}

func (r *request) decodeBody(v interface{}) (err error) {
	defer myerrors.CallWrapd(&err, r.httpRequest.Body.Close, "failed to close HTTP request body")

	decoder := json.NewDecoder(r.httpRequest.Body)
	err = decoder.Decode(v)
	if err != nil {
		return myerrors.Wrap(err, "failed to decode JSON body")
	}
	return nil
}

func (r *request) decodeQueryString(v interface{}) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag(queryTag)
	return decoder.Decode(v, r.httpRequest.URL.Query())
}

func (r *request) decodeURLParams(v interface{}) error {
	conf := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          urlTag,
		Result:           v,
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return myerrors.Wrap(err, "failed to create new decoder")
	}

	vars := mux.Vars(r.httpRequest)
	return decoder.Decode(vars)
}
