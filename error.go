package rest

type ErrorTranslator interface {
	TranslateError(error) uint
}

type ErrorTranslatorFunc func(error) uint

func (f ErrorTranslatorFunc) TranslateError(err error) uint {
	return f(err)
}

type ErrorHandler interface {
	HandleError(error)
}

type ErrorHandlerFunc func(error)

func (f ErrorHandlerFunc) HandleError(err error) {
	f(err)
}
