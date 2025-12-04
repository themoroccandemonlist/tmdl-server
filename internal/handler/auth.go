package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	session, err := h.Config.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error.", http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	session.Save(r, w)
	url := h.Config.OAuth2.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	session, _ := h.Config.Store.Get(r, "session")
	storedState, _ := session.Values["state"].(string)
	if r.URL.Query().Get("state") != storedState {
		http.Error(w, "Invalid state parameter.", http.StatusInternalServerError)
		return
	}
	delete(session.Values, "state")
	session.Save(r, w)

	code := r.URL.Query().Get("code")
	t, err := h.Config.OAuth2.Exchange(ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := h.Config.OAuth2.Client(ctx, t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var oauth2User map[string]any
	json.NewDecoder(resp.Body).Decode(&oauth2User)

	email := oauth2User["email"].(string)
	sub := oauth2User["sub"].(string)
	session.Values["user_email"] = email
	session.Values["user_sub"] = sub

	user, err := repository.GetUserByEmailAndSub(ctx, h.Config.Database, email, sub)
	var playerID *uuid.UUID
	if errors.Is(err, pgx.ErrNoRows) {
		user, _ = repository.CreateUser(ctx, h.Config.Database, email, sub)
		playerID, _ = repository.CreatePlayer(context.Background(), h.Config.Database, user.ID)
	} else {
		playerID, _ = repository.GetPlayerIDByUserID(context.Background(), h.Config.Database, user.ID)
	}

	session.Values["user_id"] = user.ID
	session.Values["user_roles"] = user.Roles
	session.Values["user_is_banned"] = user.IsBanned
	session.Values["user_is_deleted"] = user.IsDeleted
	session.Values["player_id"] = playerID
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.Config.Store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
