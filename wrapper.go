package rest

import (
	"github.com/AlexanderFadeev/myerrors"
	"net/http"
)

type Wrapper interface {
	WrapHandler(Handler) http.HandlerFunc
}

type GenericWrapper interface {
	Wrapper
	WrapGenericHandler(GenericHandler) (http.HandlerFunc, error)
	MustWrapGenericHandler(GenericHandler) http.HandlerFunc
}

type wrapper struct {
	errHandler myerrors.Handler
	translator ErrorTranslator
}

func NewWrapper(errHandler myerrors.Handler) Wrapper {
	return &wrapper{
		errHandler: errHandler,
	}
}

func NewGenericWrapper(translator ErrorTranslator, errHandler myerrors.Handler) GenericWrapper {
	return &wrapper{
		errHandler: errHandler,
		translator: translator,
	}
}

func (w *wrapper) WrapHandler(handler Handler) http.HandlerFunc {
	return WrapHandler(handler, w.errHandler)
}

func (w *wrapper) WrapGenericHandler(handler GenericHandler) (http.HandlerFunc, error) {
	return WrapGenericHandler(handler, w.translator, w.errHandler)
}

func (w *wrapper) MustWrapGenericHandler(handler GenericHandler) http.HandlerFunc {
	return MustWrapGenericHandler(handler, w.translator, w.errHandler)
}
