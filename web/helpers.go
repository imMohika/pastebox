package web

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) serverError(writer http.ResponseWriter, request *http.Request, err error) {
	s.Logger.Error(err.Error(), "method", request.Method, "uri", request.URL.RequestURI())
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *Server) clientError(writer http.ResponseWriter, _ *http.Request, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func (s *Server) render(writer http.ResponseWriter, request *http.Request, status int, page string, data templateData) {
	tmpl, ok := s.TemplateCache[page]
	if !ok {
		s.serverError(writer, request, fmt.Errorf("page %s not found in templates", page))
		return
	}

	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	writer.WriteHeader(status)
	_, err = buf.WriteTo(writer)
	if err != nil {
		s.serverError(writer, request, err)
		return
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
