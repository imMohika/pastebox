package web

import (
	"html/template"
	"net/http"
	"pastebox.mohika.ir/internal/database"
	"path/filepath"
)

type templateData struct {
	Snippet  database.Snippet
	Snippets []database.Snippet
	Form     any
	Flash    string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func (s *Server) newTemplatedata(request *http.Request) templateData {
	return templateData{
		Flash: s.SessionManager.PopString(request.Context(), "flash"),
	}
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
