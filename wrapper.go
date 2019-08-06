package rest

import "net/http"

type Wrapper interface {
	WrapHandler(Handler) http.HandlerFunc
}

type GenericWrapper interface {
	Wrapper
	WrapGenericHandler(GenericHandler) http.HandlerFunc
}

type wrapper struct {
	errChan    chan<- error
	translator ErrorTranslator
}

func NewWrapper(errChan chan<- error) Wrapper {
	return &wrapper{
		errChan: errChan,
	}
}

func NewGenericWrapper(translator ErrorTranslator, errChan chan<- error) GenericWrapper {
	return &wrapper{
		errChan:    errChan,
		translator: translator,
	}
}

func (w *wrapper) WrapHandler(handler Handler) http.HandlerFunc {
	return WrapHandler(handler, w.errChan)
}

func (w *wrapper) WrapGenericHandler(handler GenericHandler) http.HandlerFunc {
	return WrapGenericHandler(handler, w.translator, w.errChan)
}
