package rest

import "io"

type Reply interface {
	StatusCode() int
	Encode(io.Writer) error
}
