package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
	"github.com/themoroccandemonlist/tmdl-server/internal/types"
	"github.com/themoroccandemonlist/tmdl-server/internal/views/admin"
)

func (h *Handler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	const limit = 10
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	players, errr := repository.AdminGetAllPlayers(context.Background(), h.Config.Database, limit, offset, search)
	if errr != nil {
		log.Println(errr)
	}
	

	classicLevelCount, _ := repository.GetClassicLevelsCount(context.Background(), h.Config.Database)
	classicRecordCount, _ := repository.GetClassicRecordsCount(context.Background(), h.Config.Database)
	playerCount, _ := repository.GetPlayersCount(context.Background(), h.Config.Database)
	counts := types.DataCount{ClassicLevelCount: classicLevelCount, ClassicRecordCount: classicRecordCount}

	end := offset + len(players)
	start := offset + 1
	if len(players) == 0 {
		start = 0
	}
	totalPages := (playerCount + limit - 1) / limit
	pagination := types.Pagination{Start: start, End: end, Page: page, Total: playerCount, TotalPages: totalPages}
	csrfToken := csrf.TemplateField(r)
	var err error
	if IsHTMXRequest(r) {
		err = admin.Players(string(csrfToken), pagination, players).Render(r.Context(), w)
	} else {
		err = admin.Layout("Players", counts, admin.Players(string(csrfToken), pagination, players)).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}