package web

import "net/http"

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /snippet/view/{id}", s.snippetView)
	mux.HandleFunc("GET /snippet/create", s.snippetCreate)
	mux.HandleFunc("POST /snippet/create", s.snippetCreatePost)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return mux
}
