package router

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/themoroccandemonlist/tmdl-server/internal/config"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
	"github.com/themoroccandemonlist/tmdl-server/internal/middleware"
)

func New() (*mux.Router, *handler.Handler) {
	h := handler.New()

	var secure bool
	var trustedOrigins []string
	if h.Config.Environment == config.Production {
		secure = true
		trustedOrigins = append(trustedOrigins, "https://themoroccandemonlist.com", "https://www.themoroccandemonlist.com")
	} else {
		secure = false
		trustedOrigins = append(trustedOrigins, "localhost:8080", "127.0.0.1:8080")
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Use(middleware.ContentSecurityPolicy)
	r.Use(csrf.Protect(h.Config.SessionKey, csrf.Secure(secure), csrf.TrustedOrigins(trustedOrigins)))

	auth := r.PathPrefix("/").Subrouter()
	auth.Use(middleware.RequireRole(h, "USER"))
	auth.Use(middleware.RequireProfile(h))

	r.HandleFunc("/login", h.Login).Methods("GET")
	r.HandleFunc("/callback", h.Callback).Methods("GET")
	r.HandleFunc("/logout", h.Logout).Methods("GET")

	auth.HandleFunc("/profile", nil).Methods("GET")

	return r, h
}
