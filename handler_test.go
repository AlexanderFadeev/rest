package rest

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type emptyStruct struct{}

var genericHandlers = []GenericHandler{
	func(*emptyStruct) (*emptyStruct, error) { return nil, nil },
}

var badGenericHandlers = []GenericHandler{
	func() (*emptyStruct, error) { return nil, nil },
	func() (*emptyStruct, *emptyStruct) { return nil, nil },
	func(*emptyStruct) error { return nil },
	func(*emptyStruct) *emptyStruct { return nil },
	func(*emptyStruct) (*emptyStruct, *emptyStruct) { return nil, nil },
	42,
}

type mockErrorTranslator struct{}

func (mockErrorTranslator) TranslateError(err error) Error {
	return NewError(err, http.StatusInternalServerError)
}

func TestWrapGenericHandlers(t *testing.T) {
	errChan := make(chan error)

	for _, gh := range genericHandlers {
		MustWrapGenericHandler(gh, new(mockErrorTranslator), errChan)
	}
}

func TestWrapBadGenericHandlers(t *testing.T) {
	errChan := make(chan error)

	for index, gh := range badGenericHandlers {
		func() {
			defer func() {
				recover()
			}()

			MustWrapGenericHandler(gh, new(mockErrorTranslator), errChan)
			assert.Failf(t, "MustWrapGenerichandler should panic", "handler index %d", index)
		}()
	}
}

func TestCallGenericWrapper(t *testing.T) {
	called := false

	type args struct {
		X int
	}
	type reply struct {
		X int
	}

	handler := func(a *args) (*reply, error) {
		called = true
		assert.Equal(t, 42, a.X)
		return &reply{a.X}, nil
	}

	wrapped := MustWrapGenericHandler(handler, mockErrorTranslator{}, nil)

	body := []byte(`{"X": 42}`)
	bodyBuf := bytes.NewBuffer(body)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/test", bodyBuf)
	assert.Nil(t, err)

	wrapped.ServeHTTP(w, r)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)

	var replyMap map[string]interface{}
	d := json.NewDecoder(w.Body)
	err = d.Decode(&replyMap)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"X": 42.}, replyMap)
}
