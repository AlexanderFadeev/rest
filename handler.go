package rest

import (
	"bytes"
	"io"
	"net/http"

	"github.com/afadeevz/omnierrors"
	"github.com/mailru/easyjson"
)

type Handler[Req any, Resp any, ReqPtr PtrToRequest[Req], RespPtr PtrToResponse[Resp]] func(*Req) (*Resp, error)

func (h Handler[Req, Rep, ReqPtr, RespPtr]) ToHTTPHandler(errTranslator ErrorTranslator, errHandler ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := func() (*bytes.Buffer, error) {
			var req ReqPtr = new(Req)
			err := decodeRequest(r, req)
			if err != nil {
				return nil, omnierrors.Wrap(err, "failed to decode REST request")
			}

			resp, err := h(req)
			if err != nil {
				return nil, omnierrors.Wrap(err, "failed to handle REST request")
			}

			var buf bytes.Buffer
			_, err = easyjson.MarshalToWriter(RespPtr(resp), &buf)
			if err != nil {
				return nil, omnierrors.Wrap(err, "failed to encode REST response")
			}

			return &buf, nil
		}()

		if err != nil {
			errHandler.HandleError(err)

			status := errTranslator.TranslateError(err)
			w.WriteHeader(int(status))

			errResp := newErrorResponse(err)
			_, err = easyjson.MarshalToWriter(errResp, w)
			if err != nil {
				errHandler.HandleError(err)
			}

			return
		}

		_, err = io.Copy(w, buf)
		if err != nil {
			err = omnierrors.Wrap(err, "failed to send REST response")
			errHandler.HandleError(err)
		}
	}
}

func Wrap[Req any, Resp any, ReqPtr PtrToRequest[Req], RespPtr PtrToResponse[Resp]](handler Handler[Req, Resp, ReqPtr, RespPtr], errTranslator ErrorTranslator, errHandler ErrorHandler) http.HandlerFunc {
	return handler.ToHTTPHandler(errTranslator, errHandler)
}
