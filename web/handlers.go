package web

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pastebox.mohika.ir/internal/database"
	"pastebox.mohika.ir/internal/validator"
	"strconv"
	"time"
)

func (s *Server) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := s.Queries.GetLatestSnippets(s.Ctx)
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	data := s.newTemplatedata(request)
	data.Snippets = snippets

	s.render(writer, request, http.StatusOK, "home.gohtml", data)
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

	data := s.newTemplatedata(request)
	data.Snippet = snippet

	s.render(writer, request, http.StatusOK, "snippet_view.gohtml", data)
}

func (s *Server) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	data := s.newTemplatedata(request)
	data.Form = snippetCreateForm{
		Expires: -1,
	}

	s.render(writer, request, http.StatusOK, "snippet_create.gohtml", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (s *Server) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		s.clientError(writer, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		s.clientError(writer, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   request.PostForm.Get("title"),
		Content: request.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := s.newTemplatedata(request)
		data.Form = form
		s.render(writer, request, http.StatusUnprocessableEntity, "snippet_create.gohtml", data)
		return
	}

	var expiresTime sql.NullTime
	if form.Expires <= 0 {
		expiresTime = sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	} else {
		expiresTime = sql.NullTime{
			Time:  time.Now().Add(time.Duration(expires) * time.Hour * 24),
			Valid: true,
		}
	}

	id, err := s.Queries.CreateSnippet(s.Ctx, database.CreateSnippetParams{
		Title:   form.Title,
		Content: form.Content,
		Expires: expiresTime,
	})
	if err != nil {
		s.serverError(writer, request, err)
		return
	}

	s.SessionManager.Put(request.Context(), "flash", "Snippet successfully created")

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
