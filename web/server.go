package web

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log/slog"
	"pastebox.mohika.ir/internal/database"
)

type Server struct {
	Logger         *slog.Logger
	Queries        *database.Queries
	Ctx            context.Context
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
}
