package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/themoroccandemonlist/tmdl-server/internal/config"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type Handler struct {
	Config *config.Config
	Validate *validator.Validate
}

func New() *Handler {
	cfg := config.New()

	v := validator.New()

	v.RegisterValidation("difficulty", func(fl validator.FieldLevel) bool {
		return enum.Difficulty(fl.Field().String()).IsValid()
	})

	v.RegisterValidation("duration", func(fl validator.FieldLevel) bool {
		return enum.Duration(fl.Field().String()).IsValid()
	})

	return &Handler{
		Config: cfg,
		Validate: v,
	}
}

func IsHTMXRequest(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}

