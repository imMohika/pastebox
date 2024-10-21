package web

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pastebox.mohika.ir/internal/database"
	"strconv"
)

func (s *Server) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := s.Queries.GetLatestSnippets(s.Ctx)
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	s.render(writer, request, http.StatusOK, "home.gohtml", templateData{
		Snippets: snippets,
	})
}

func (s *Server) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	snippet, err := s.Queries.GetSnippet(s.Ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(writer, request)
			return
		}
		s.serverError(writer, request, err)
		return
	}

	s.render(writer, request, http.StatusOK, "home.gohtml", templateData{
		Snippet: snippet,
	})
}

func (s *Server) snippetCreate(writer http.ResponseWriter, _ *http.Request) {
	writer.Write([]byte("form"))
}

func (s *Server) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	title := "Meow"
	content := "meow\nmeow\nmeow"

	id, err := s.Queries.CreateSnippet(s.Ctx, database.CreateSnippetParams{
		Title:   title,
		Content: content,
	})
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	http.RedirectHandler(fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
