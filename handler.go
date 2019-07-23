package rest

import (
	"fmt"
	"net/http"
)

type Handler func(Request) Reply

type handlerWrapper struct {
	handler   Handler
	errorChan chan<- error
}

func (hw *handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := hw.serveHTTPImpl(r, w)
	if err != nil {
		hw.errorChan <- fmt.Errorf("failed to handle HTTP request: %w", err)
	}
}

func (hw *handlerWrapper) serveHTTPImpl(r *http.Request, w http.ResponseWriter) error {
	request := request{httpRequest: r}
	reply := hw.handler(&request)
	err := hw.encodeReply(w, reply)
	return err
}

func (hw *handlerWrapper) encodeReply(w http.ResponseWriter, reply Reply) (err error) {
	w.WriteHeader(reply.StatusCode())

	err = reply.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode reply")
	}

	return
}
