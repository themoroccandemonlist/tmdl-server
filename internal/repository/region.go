package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

// to verify
func UpdateRegionsClassicPoints(ctx context.Context, pool *pgxpool.Pool, regions []*model.Region) error {
    ids := make([]uuid.UUID, len(regions))
    points := make([]decimal.Decimal, len(regions))
    for i, r := range regions {
        ids[i] = r.ID
        points[i] = r.ClassicPoints
    }

    _, err := pool.Exec(ctx, `
        UPDATE regions SET classic_points = data.points
        FROM unnest($1::uuid[], $2::numeric[]) AS data(id, points)
        WHERE regions.id = data.id
    `, ids, points)
    if err != nil {
        return fmt.Errorf("failed to bulk update region points: %w", err)
    }
    return nil
}