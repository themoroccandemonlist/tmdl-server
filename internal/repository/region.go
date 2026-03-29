package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/dto"
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

func AdminGetAllRegions(ctx context.Context, pool *pgxpool.Pool) ([]dto.AdminRegionRow, error) {
	query := `
		SELECT r.id, r.name, COUNT(p.id) as total_players, r.classic_points, r.platformer_points
		FROM regions r
		LEFT JOIN players p ON p.region_id = r.id
		GROUP BY r.id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query regions: %w", err)
	}
	defer rows.Close()

	var regions []dto.AdminRegionRow
	for rows.Next() {
		var row dto.AdminRegionRow
		if err := rows.Scan(&row.ID, &row.Name, &row.TotalPlayers, &row.ClassicPoints, &row.PlatformerPoints); err != nil {
			return nil, fmt.Errorf("failed to scan region row: %w", err)
		}
		regions = append(regions, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
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