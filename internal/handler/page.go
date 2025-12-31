package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
	"github.com/themoroccandemonlist/tmdl-server/internal/views"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := views.Home().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ProfileSetup(w http.ResponseWriter, r *http.Request) {
	regions, err := repository.GetAllRegionIDsAndNames(context.Background(), h.Config.Database)
	if err != nil {
		http.Error(w, "Failed to load regions", http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.TemplateField(r)

	err = views.ProfileSetup(string(csrfToken), regions).Render(context.Background(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ProfileSetupSubmit(w http.ResponseWriter, r *http.Request) {
	session, _ := h.Config.Store.Get(r, "session")
	playerID, ok := session.Values["player_id"].(uuid.UUID)
	if !ok {
		http.Error(w, "Player ID not found", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := h.Config.Sanitizer.Sanitize(r.FormValue("username"))
	regionID, _ := uuid.Parse(h.Config.Sanitizer.Sanitize(r.FormValue("region")))

	err := repository.UpdateUsernameAndRegion(context.Background(), h.Config.Database, playerID, username, regionID)
	if err != nil {
		log.Printf("Failed to update player profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	session.Values["player_username"] = username
	session.Values["player_region_id"] = regionID
	session.Save(r, w)

	w.Header().Set("HX-Redirect", "/profile")
	w.WriteHeader(http.StatusOK)
}
