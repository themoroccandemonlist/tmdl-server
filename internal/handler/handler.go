package handler

import (
	"html/template"

	"github.com/themoroccandemonlist/tmdl-server/internal/config"
	"github.com/themoroccandemonlist/tmdl-server/templates"
)

type Handler struct {
	Config *config.Config
}

func New() *Handler {
	cfg := config.New()
	return &Handler{
		Config: cfg,
	}
}

var baseTemplate = template.Must(template.ParseFS(templates.TemplatesFS, "layout.html"))

func composeTemplate(templateNames ...string) *template.Template {
	tmpl := template.Must(baseTemplate.Clone())

	paths := []string{}
	for _, name := range templateNames {
		paths = append(paths, name+".html")
	}
	return template.Must(tmpl.ParseFS(templates.TemplatesFS, paths...))
}
