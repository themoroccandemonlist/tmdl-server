package router

import (
	"github.com/gorilla/mux"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func New() (*mux.Router, *handler.Handler) {
	h := handler.New()
	r := mux.NewRouter()

	return r, h
}
