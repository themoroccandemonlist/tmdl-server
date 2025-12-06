package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func GetAllRegionIDsAndNames(ctx context.Context, pool *pgxpool.Pool) ([]*model.Region, error) {
	query := `
		SELECT id, name
		FROM regions
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Printf("Unable to fetch resources: %v", err)
		return nil, err
	}
	defer rows.Close()

	var regions []*model.Region
	for rows.Next() {
		var region model.Region
		rows.Scan(&region.ID, &region.Name)
		regions = append(regions, &region)
	}
	return regions, nil
}
