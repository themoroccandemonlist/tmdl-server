package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
)

func (h *Handler) ProfileSetup(w http.ResponseWriter, r *http.Request) {
	regions, err := repository.GetAllRegionIDsAndNames(context.Background(), h.Config.Database)
	if err != nil {
		http.Error(w, "Failed to load regions", http.StatusInternalServerError)
		return
	}

	data := struct {
		CSRFField template.HTML
		Regions   []*model.Region
	}{
		CSRFField: csrf.TemplateField(r),
		Regions:   regions,
	}

	err = TMPL_PROFILE_SETUP.ExecuteTemplate(w, "layout", data)
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

	username := r.FormValue("username")
	regionID, _ := uuid.Parse(r.FormValue("region"))

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
