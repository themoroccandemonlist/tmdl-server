package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
	"github.com/themoroccandemonlist/tmdl-server/internal/repository"
	"github.com/themoroccandemonlist/tmdl-server/internal/types"
	"github.com/themoroccandemonlist/tmdl-server/internal/util"
	"github.com/themoroccandemonlist/tmdl-server/internal/views/admin"
)

func (h *Handler) ListClassicLevels(w http.ResponseWriter, r *http.Request) {
	const limit = 10
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	classicLevels, _ := repository.AdminGetAllClassicLevels(context.Background(), h.Config.Database, limit, offset, search)
	classicLevelCount, _ := repository.GetClassicLevelsCount(context.Background(), h.Config.Database)
	counts := types.DataCount{ClassicLevelCount: classicLevelCount}
	end := offset + len(classicLevels)
	start := offset + 1
	if len(classicLevels) == 0 {
		start = 0
	}
	totalPages := (classicLevelCount + limit - 1) / limit
	pagination := types.Pagination{Start: start, End: end, Page: page, Total: classicLevelCount, TotalPages: totalPages}
	csrfToken := csrf.TemplateField(r)
	var err error
	if IsHTMXRequest(r) {
		err = admin.ClassicLevels(string(csrfToken), pagination, classicLevels).Render(r.Context(), w)
	} else {
		err = admin.Layout("Admin - Classic Levels", counts, admin.ClassicLevels(string(csrfToken), pagination, classicLevels)).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateClassicLevel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, `failed to parse form`, http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			http.Error(w, `thumbnail is required`, http.StatusBadRequest)
		} else {
			http.Error(w, `failed to read thumbnail`, http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	if !util.IsValidImageType(header.Header.Get("Content-Type")) {
		http.Error(w, `thumbnail must be an image (jpeg, png, gif)`, http.StatusBadRequest)
		return
	}

	thumbnailPath, err := util.SaveThumbnail(file, header)
	if err != nil {
		http.Error(w, `failed to save thumbnail`, http.StatusInternalServerError)
		return
	}

	ranking, err := strconv.Atoi(r.FormValue("ranking"))
	if err != nil {
		http.Error(w, `ranking must be a number`, http.StatusBadRequest)
		return
	}

	listPercentage, err := strconv.Atoi(r.FormValue("list_percentage"))
	if err != nil {
		http.Error(w, `list_percentage must be a number`, http.StatusBadRequest)
		return
	}

	level := model.ClassicLevel{
		LevelID:        h.Config.Sanitizer.Sanitize(r.FormValue("level_id")),
		Name:           h.Config.Sanitizer.Sanitize(r.FormValue("name")),
		Publisher:      h.Config.Sanitizer.Sanitize(r.FormValue("publisher")),
		Difficulty:     enum.Difficulty(h.Config.Sanitizer.Sanitize(r.FormValue("difficulty"))),
		Duration:       enum.Duration(h.Config.Sanitizer.Sanitize(r.FormValue("duration"))),
		Ranking:        ranking,
		ListPercentage: listPercentage,
		YoutubeLink:    h.Config.Sanitizer.Sanitize(r.FormValue("youtube_link")),
		ThumbnailPath:  "/" + thumbnailPath,
	}

	if err := h.Validate.Struct(level); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, validationErrors[0].Field()+` is invalid or missing`, http.StatusBadRequest)
			return
		}
		http.Error(w, `validation failed`, http.StatusBadRequest)
		return
	}

	level.CalculatePoints()

	if err := repository.CreateClassicLevel(context.Background(), h.Config.Database, level); err != nil {
		log.Printf("failed to insert classic level: %v", err)
		http.Error(w, "failed to add level", http.StatusInternalServerError)
		return
	}

	affectedLevels, err := repository.GetAffectedClassicLevels(context.Background(), h.Config.Database, level.Ranking)
	if err != nil {
		log.Printf("failed to fetch affected levels: %v", err)
		http.Error(w, "failed to fetch affected levels", http.StatusInternalServerError)
	}

	var records []*model.ClassicRecord

	if len(affectedLevels) <= 1 {
		w.WriteHeader(http.StatusCreated)
	}

	for _, affectedLevel := range affectedLevels {
		if !(affectedLevel.ID == level.ID) {
			affectedLevel.Ranking = affectedLevel.Ranking + 1
		}
		if affectedLevel.Ranking > 75 {
			affectedLevel.ListPercentage = 100
		}
		if affectedLevel.Ranking > 150 {
			affectedLevel.OldPoints = affectedLevel.Points
			affectedLevel.OldMinimumPoints = affectedLevel.MinimumPoints
			affectedLevel.Points = decimal.Zero
			affectedLevel.MinimumPoints = decimal.Zero
			affectedRecords, err := repository.GetAllClassicRecordsByClassicLevel(context.Background(), h.Config.Database, affectedLevel.ID)
			if err != nil {
				log.Printf("failed to fetch affected records: %v", err)
				http.Error(w, "failed to fetch affected records", http.StatusInternalServerError)
			}
			records = append(records, affectedRecords...)
			continue
		}
		if affectedLevel.ID == level.ID {
			continue
		}
		affectedLevel.OldPoints = affectedLevel.Points
		affectedLevel.OldMinimumPoints = affectedLevel.MinimumPoints
		affectedLevel.CalculatePoints()
		affectedRecords, err := repository.GetAllClassicRecordsByClassicLevel(context.Background(), h.Config.Database, affectedLevel.ID)
		if err != nil {
			log.Printf("failed to fetch affected records: %v", err)
			http.Error(w, "failed to fetch affected records", http.StatusInternalServerError)
		}
		records = append(records, affectedRecords...)
	}
	if len(records) == 0 {
		w.WriteHeader(http.StatusCreated)
	}

	var updatedPlayers []*model.Player
	var updatedRegions []*model.Region

	for _, record := range records {
		var level *model.ClassicLevel
		for _, al := range affectedLevels {
			if al.ID == record.ClassicLevelID {
				level = al
				break
			}
		}
		if level == nil {
			continue
		}

		player := record.Player
		region := record.Player.Region

		var oldValue, newValue decimal.Decimal

		if record.Progress == 100 {
			oldValue = level.OldPoints
			newValue = level.Points
		} else {
			oldValue = level.OldMinimumPoints
			newValue = level.MinimumPoints
		}

		oldPlayerPoints := player.ClassicPoints
		newPlayerPoints := oldPlayerPoints.Sub(oldValue).Add(newValue)

		oldRegionPoints := region.ClassicPoints
		newRegionPoints := oldRegionPoints.Sub(oldValue).Add(newValue)

		if !oldPlayerPoints.Equal(newPlayerPoints) {
			player.ClassicPoints = newPlayerPoints
			updatedPlayers = append(updatedPlayers, player)
		}
		if !oldRegionPoints.Equal(newRegionPoints) {
			region.ClassicPoints = newRegionPoints
			updatedRegions = append(updatedRegions, region)
		}
	}

	if len(updatedPlayers) > 0 {
		if err := repository.UpdatePlayersClassicPoints(r.Context(), h.Config.Database, updatedPlayers); err != nil {
			http.Error(w, `failed to update player points`, http.StatusInternalServerError)
			return
		}
	}

	if len(updatedRegions) > 0 {
		if err := repository.UpdateRegionsClassicPoints(r.Context(), h.Config.Database, updatedRegions); err != nil {
			http.Error(w, `failed to update region points`, http.StatusInternalServerError)
			return
		}
	}

	if len(affectedLevels) > 0 {
		if err := repository.UpdateClassicLevelsPoints(r.Context(), h.Config.Database, affectedLevels); err != nil {
			log.Printf("failed to bulk update classic levels: %v", err)
			http.Error(w, `failed to update affected levels`, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
