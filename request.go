package rest

import (
	"net/http"
	"net/url"

	"github.com/afadeevz/omnierrors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mailru/easyjson"
	"github.com/mitchellh/mapstructure"
)

const (
	queryTag = "query"
	urlTag   = "url"
)

type Request interface {
	easyjson.Unmarshaler
}

type PtrToRequest[T any] interface {
	Request
	*T
}

func decodeRequest(httpReq *http.Request, req Request) error {
	err := decodeBody(httpReq, req)
	if err != nil {
		return omnierrors.Wrap(err, "failed to decode request body")
	}

	err = decodeQueryString(httpReq.URL.Query(), req)
	if err != nil {
		return omnierrors.Wrap(err, "failed to decode query strings")
	}

	vars := mux.Vars(httpReq)
	err = decodeURLParams(vars, req)
	return omnierrors.Wrap(err, "failed to decode URL params")
}

func decodeBody(httpReq *http.Request, req Request) error {
	if httpReq.ContentLength == 0 {
		return nil
	}

	err := easyjson.UnmarshalFromReader(httpReq.Body, req)
	return omnierrors.Wrap(err, "failed to decode JSON body")
}

func decodeQueryString(query url.Values, req Request) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag(queryTag)
	return decoder.Decode(req, query)
}

func decodeURLParams(vars map[string]string, req Request) error {
	conf := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          urlTag,
		Result:           req,
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return omnierrors.Wrap(err, "failed to create new decoder")
	}

	return decoder.Decode(vars)
}
