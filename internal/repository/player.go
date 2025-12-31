package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func GetPlayerByUserID(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*model.Player, error) {
	query := `
		SELECT id, username, avatar, region_id
		FROM players
		WHERE user_id = $1
	`

	var player model.Player
	err := pool.QueryRow(ctx, query, userID).Scan(&player.ID, &player.Username, &player.Avatar, &player.RegionID)
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
