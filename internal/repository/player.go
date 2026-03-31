package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/dto"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func AdminGetAllPlayers(ctx context.Context, pool *pgxpool.Pool, limit int, offset int, search string) ([]dto.AdminPlayerRow, error) {
	// query := `
	// 	SELECT p.id, p.username, r.id, r.name, p.classic_points, p.platformer_points,
	// 	       COUNT(cr.id) AS total_classic_records, COUNT(pr.id) AS total_platformer_records, p.status
	// 	FROM players p
	// 	LEFT JOIN regions r ON r.id = p.region_id
	// 	LEFT JOIN classic_records cr ON cr.player_id = p.id
	// 	LEFT JOIN platformer_records pr ON pr.player_id = p.id
	// 	WHERE ($3 = '' OR p.username ILIKE '%' || $3 || '%')
	// 	GROUP BY p.id, p.username, r.id, r.name, p.classic_points, p.platformer_points, p.status
	// 	ORDER BY p.username ASC
	// 	LIMIT $1 OFFSET $2
	// `
	query := `
		SELECT p.id, p.username, r.id, r.name, p.classic_points, p.platformer_points,
		       COUNT(cr.id) AS total_classic_records, p.status
		FROM players p
		LEFT JOIN regions r ON r.id = p.region_id
		LEFT JOIN classic_records cr ON cr.player_id = p.id
		WHERE ($3 = '' OR p.username ILIKE '%' || $3 || '%')
		GROUP BY p.id, p.username, r.id, r.name, p.classic_points, p.platformer_points, p.status
		ORDER BY p.username ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := pool.Query(ctx, query, limit, offset, search)
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %w", err)
	}
	defer rows.Close()

	var players []dto.AdminPlayerRow
	for rows.Next() {
		var row dto.AdminPlayerRow
		// if err := rows.Scan(&row.ID, &row.Username, &row.Region.ID, &row.Region.Name, &row.ClassicPoints, &row.PlatformerPoints, &row.TotalClassicRecords, &row.TotalPlatformerRecords, &row.Status); err != nil {
		if err := rows.Scan(&row.ID, &row.Username, &row.Region.ID, &row.Region.Name, &row.ClassicPoints, &row.PlatformerPoints, &row.TotalClassicRecords, &row.Status); err != nil {
			return nil, fmt.Errorf("failed to scan player row: %w", err)
		}
		players = append(players, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return players, nil
}

func GetPlayersCount(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM players
	`

	var count int
	if err := pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count players: %w", err)
	}
	return count, nil
}

func GetPlayerByUserID(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*model.Player, error) {
	query := `
		SELECT id, username, avatar, region_id, status
		FROM players
		WHERE user_id = $1
	`

	var player model.Player
	err := pool.QueryRow(ctx, query, userID).Scan(&player.ID, &player.Username, &player.Avatar, &player.RegionID, &player.Status)
	if err != nil {
		log.Printf("Unable to fetch resource: %v", err)
		return nil, err
	}
	return &player, nil
}

func CreatePlayer(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*model.Player, error) {
	query := `
		INSERT INTO players (user_id, status)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id, username, region_id, status
	`

	var player model.Player
	err := pool.QueryRow(ctx, query, userID, enum.Setup).Scan(&player.ID, &player.Username, &player.RegionID, &player.Status)
	if err != nil {
		log.Printf("Unable to insert resource: %v", err)
		return nil, err
	}
	return &player, nil
}

func UpdateUsernameAndRegion(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID, username string, regionID uuid.UUID) error {
	query := `
		UPDATE players
		SET username = $2, region_id = $3
		WHERE id = $1
	`

	commandTag, err := pool.Exec(ctx, query, playerID, username, regionID)
	if err != nil {
		log.Printf("Unable to insert resource: %v", err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

// to verify
func UpdatePlayersClassicPoints(ctx context.Context, pool *pgxpool.Pool, players []*model.Player) error {
    ids := make([]uuid.UUID, len(players))
    points := make([]decimal.Decimal, len(players))
    for i, p := range players {
        ids[i] = p.ID
        points[i] = p.ClassicPoints
    }

    _, err := pool.Exec(ctx, `
        UPDATE players SET classic_points = data.points
        FROM unnest($1::uuid[], $2::numeric[]) AS data(id, points)
        WHERE players.id = data.id
    `, ids, points)
    if err != nil {
        return fmt.Errorf("failed to bulk update player points: %w", err)
    }
    return nil
}