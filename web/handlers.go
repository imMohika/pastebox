package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (s *Server) home(writer http.ResponseWriter, request *http.Request) {
	files := []string{
		"ui/html/base.gohtml",
		"ui/html/partials/nav.gohtml",
		"ui/html/pages/home.gohtml",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	err = t.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		s.serverError(writer, request, err)
	}
}

func (s *Server) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "id: %d", id)
}

func (s *Server) snippetCreate(writer http.ResponseWriter, _ *http.Request) {
	writer.Write([]byte("form"))
}

func (s *Server) snippetCreatePost(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("save"))
}
