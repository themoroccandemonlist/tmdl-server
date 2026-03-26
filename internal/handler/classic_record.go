package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
	"github.com/themoroccandemonlist/tmdl-server/internal/types"
	"github.com/themoroccandemonlist/tmdl-server/internal/views/admin"
)

func (h *Handler) ListClassicRecords(w http.ResponseWriter, r *http.Request) {
	const limit = 10
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	classicRecords, _ := repository.AdminGetAllClassicRecords(context.Background(), h.Config.Database, limit, offset, search)

	classicLevelCount, _ := repository.GetClassicLevelsCount(context.Background(), h.Config.Database)
	classicRecordCount, _ := repository.GetClassicRecordsCount(context.Background(), h.Config.Database)
	counts := types.DataCount{ClassicLevelCount: classicLevelCount, ClassicRecordCount: classicRecordCount}

	end := offset + len(classicRecords)
	start := offset + 1
	if len(classicRecords) == 0 {
		start = 0
	}
	totalPages := (classicRecordCount + limit - 1) / limit
	pagination := types.Pagination{Start: start, End: end, Page: page, Total: classicRecordCount, TotalPages: totalPages}
	csrfToken := csrf.TemplateField(r)
	var err error
	if IsHTMXRequest(r) {
		err = admin.ClassicRecords(string(csrfToken), pagination, classicRecords).Render(r.Context(), w)
	} else {
		err = admin.Layout("Admin - Classic Records", counts, admin.ClassicRecords(string(csrfToken), pagination, classicRecords)).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}