package rest

type Route struct {
	Pattern string
	Handler Handler
	Method  string
}

type Config struct {
	Routes []Route
}
