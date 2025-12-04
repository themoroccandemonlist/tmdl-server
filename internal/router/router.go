package router

import (
	"github.com/gorilla/mux"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
	"github.com/themoroccandemonlist/tmdl-server/internal/middleware"
)

func New() (*mux.Router, *handler.Handler) {
	h := handler.New()
	r := mux.NewRouter()
	r.Use(middleware.ContentSecurityPolicy)

	return r, h
}
