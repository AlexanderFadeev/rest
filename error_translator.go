package rest

type ErrorTranslator interface {
	TranslateError(error) Error
}
