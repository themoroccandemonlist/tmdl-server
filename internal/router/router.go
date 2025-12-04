package router

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/themoroccandemonlist/tmdl-server/internal/config"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
	"github.com/themoroccandemonlist/tmdl-server/internal/middleware"
)

func New() (*mux.Router, *handler.Handler) {
	h := handler.New()

	var secure bool
	if h.Config.Environment == config.Production {
		secure = true
	} else {
		secure = false
	}

	r := mux.NewRouter()

	r.Use(middleware.ContentSecurityPolicy)
	r.Use(csrf.Protect(h.Config.SessionKey, csrf.Secure(secure)))

	return r, h
}
