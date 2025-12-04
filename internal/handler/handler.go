package handler

import "github.com/themoroccandemonlist/tmdl-server/internal/config"

type Handler struct {
	Config *config.Config
}

func New() *Handler {
	cfg := config.New()
	return &Handler{
		Config: cfg,
	}
}
