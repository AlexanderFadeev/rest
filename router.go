package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	http.Handler
	ErrorChan() <-chan error
}

type router struct {
	*mux.Router
	errorChan chan error
}

func NewRouter(config *Config) Router {
	router := newRouter()

	for _, route := range config.Routes {
		router.handle(route.Pattern, route.Handler, route.Method)
	}

	return router
}

func newRouter() *router {
	return &router{
		Router:    mux.NewRouter(),
		errorChan: make(chan error),
	}
}

func (r *router) handle(pattern string, handler Handler, method string) {
	r.Handle(pattern, &handlerWrapper{
		handler:   handler,
		errorChan: r.errorChan,
	}).Methods(method)
}

func (r *router) ErrorChan() <-chan error {
	return r.errorChan
}
