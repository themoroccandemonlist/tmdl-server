package handler

import (
	"context"
	"net/http"

	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
	"github.com/themoroccandemonlist/tmdl-server/internal/types"
	"github.com/themoroccandemonlist/tmdl-server/internal/views/admin"
)

func (h *Handler) ListRegions(w http.ResponseWriter, r *http.Request) {
	regions, _ := repository.AdminGetAllRegions(context.Background(), h.Config.Database)

	classicLevelCount, _ := repository.GetClassicLevelsCount(context.Background(), h.Config.Database)
	classicRecordCount, _ := repository.GetClassicRecordsCount(context.Background(), h.Config.Database)
	counts := types.DataCount{ClassicLevelCount: classicLevelCount, ClassicRecordCount: classicRecordCount, }

	var err error
	if IsHTMXRequest(r) {
		err = admin.Regions(regions).Render(r.Context(), w)
	} else {
		err = admin.Layout("Regions", counts, admin.Regions(regions)).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}