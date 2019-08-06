package rest

import (
	"bytes"
	"fmt"
	"github.com/AlexanderFadeev/myerrors"
	"io"
	"net/http"
	"reflect"
)

type (
	GenericHandler interface{} // TODO: Go2: func(type ArgsType, ReplyType)(*ArgsType) (*ReplyType, error)
	Handler        func(Request) Reply
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func WrapGenericHandler(genericHandler GenericHandler, translator ErrorTranslator, errChan chan<- error) http.HandlerFunc {
	rv := reflect.ValueOf(genericHandler)
	rt := rv.Type()
	if rk := rt.Kind(); rk != reflect.Func {
		panic(myerrors.Errorf("expected function, got %s", rk.String()))
	}

	if numIn := rt.NumIn(); numIn != 1 {
		panic(myerrors.Errorf("expected function with 1 argument, got %d", numIn))
	}

	if numOut := rt.NumOut(); numOut != 2 {
		panic(myerrors.Errorf("expected function with 2 return values, got %d)", numOut))
	}

	if !rt.Out(1).Implements(errorInterface) {
		panic(myerrors.Errorf("expected second return value to implement error interface, got %s", rt.Out(1).Name()))
	}

	handler := func(req Request) Reply {
		in := reflect.New(rt.In(0).Elem())
		err := req.Decode(in.Interface())
		if err != nil {
			return NewError(err, http.StatusBadRequest)
		}

		ret := rv.Call([]reflect.Value{in})
		iErr := ret[1].Interface()
		if iErr != nil {
			return translator.TranslateError(iErr.(error))
		}

		out := ret[0].Elem().Interface()
		return NewOKReply(out)
	}

	return WrapHandler(handler, errChan)
}

func WrapHandler(handler Handler, errChan chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := request{httpRequest: r}
		reply := handler(&request)
		err := encodeReply(w, reply)
		if err != nil {
			errChan <- myerrors.Wrap(err, "failed to handle HTTP request")
		}
	}
}

func encodeReply(w http.ResponseWriter, reply Reply) error {
	var buf bytes.Buffer

	err := reply.Encode(&buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to encode reply")
	}

	w.WriteHeader(reply.StatusCode())
	_, err = io.Copy(w, &buf)
	if err != nil {
		return myerrors.Wrap(err, "failed to write reply")
	}

	return nil
}
