package web

import (
	"github.com/justinas/alice"
	"net/http"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(s.SessionManager.LoadAndSave)
	standard := alice.New(s.panicRecover, s.logRequest, commonHeaders)

	mux.Handle("GET /{$}", dynamic.ThenFunc(s.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(s.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(s.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(s.snippetCreatePost))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return standard.Then(mux)
}
