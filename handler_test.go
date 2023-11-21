package rest_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afadeevz/omnierrors"
	"github.com/afadeevz/rest"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

//easyjson:json
type SetReq struct {
	Key   int `json:"key"`
	Value int `json:"value"`
}

//easyjson:json
type SetResp struct{}

//easyjson:json
type GetReq struct {
	Key int `json:"key"`
}

//easyjson:json
type GetResp struct {
	Value int `json:"value"`
}

//easyjson:json
type ErrorResp struct {
	Error string `json:"error"`
}

type stubHandler struct {
	data map[int]int
}

func newStubHandler() *stubHandler {
	return &stubHandler{
		data: make(map[int]int),
	}
}

var (
	errNotFound      = omnierrors.New("not found")
	errAlreadyExists = omnierrors.New("already exists")
)

func (mh *stubHandler) Set(req *SetReq) (*SetResp, error) {
	if _, ok := mh.data[req.Key]; ok {
		return nil, errAlreadyExists
	}

	mh.data[req.Key] = req.Value
	return &SetResp{}, nil
}

func (mh *stubHandler) Get(req *GetReq) (*GetResp, error) {
	if _, ok := mh.data[req.Key]; !ok {
		return nil, errNotFound
	}

	return &GetResp{
		Value: mh.data[req.Key],
	}, nil
}

func translateError(err error) uint {
	switch {
	case omnierrors.Is(err, errNotFound):
		return http.StatusNotFound
	case omnierrors.Is(err, errAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func handleError(error) {
	// ignore
}

func TestCallGenericWrapper(t *testing.T) {
	handler := newStubHandler()
	et := rest.ErrorTranslatorFunc(translateError)
	eh := rest.ErrorHandlerFunc(handleError)

	set := rest.Wrap(handler.Set, et, eh)
	get := rest.Wrap(handler.Get, et, eh)

	k := 42
	v := 69

	checkHandler(t, get, GetReq{Key: k}, ErrorResp{"failed to handle REST request: not found"}, http.StatusNotFound)
	checkHandler(t, set, SetReq{Key: k, Value: v}, SetResp{}, http.StatusOK)
	checkHandler(t, set, SetReq{Key: k, Value: v}, ErrorResp{"failed to handle REST request: already exists"}, http.StatusConflict)
	checkHandler(t, get, GetReq{Key: k}, GetResp{Value: v}, http.StatusOK)
}

func checkHandler(t *testing.T, handler http.HandlerFunc, req easyjson.Marshaler, resp easyjson.Marshaler, status int) {
	var bodyBuf bytes.Buffer
	_, err := easyjson.MarshalToWriter(req, &bodyBuf)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/test", &bodyBuf)
	assert.Nil(t, err)

	handler(w, r)
	assert.Equal(t, status, w.Code)

	var respBuf bytes.Buffer
	_, err = easyjson.MarshalToWriter(resp, &respBuf)
	assert.Nil(t, err)

	assert.Equal(t, respBuf.String(), w.Body.String())
}
