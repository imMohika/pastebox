package main

import (
	"net/http"
	"runtime/debug"
)

func (s *Server) serverError(writer http.ResponseWriter, request *http.Request, err error) {
	s.logger.Error(err.Error(), "method", request.Method, "uri", request.URL.RequestURI(), "trace", string(debug.Stack()))
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *Server) clientError(writer http.ResponseWriter, _ *http.Request, status int) {
	http.Error(writer, http.StatusText(status), status)
}
