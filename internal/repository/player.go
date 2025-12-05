package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func GetPlayerIDAndNameByUserID(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*model.Player, error) {
	query := `
		SELECT id, username, region_id
		FROM players
		WHERE user_id = $1
	`

	var player model.Player
	err := pool.QueryRow(ctx, query, userID).Scan(&player.ID, &player.Username, &player.RegionID)
	if err != nil {
		log.Printf("Unable to fetch resource: %v", err)
		return nil, err
	}
	return &player, nil
}

func CreatePlayer(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*model.Player, error) {
	query := `
		INSERT INTO players (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id, username, region_id
	`

	var player model.Player
	err := pool.QueryRow(ctx, query, userID).Scan(&player.ID, &player.Username, &player.RegionID)
	if err != nil {
		log.Printf("Unable to insert resource: %v", err)
		return nil, err
	}
	return &player, nil
}
