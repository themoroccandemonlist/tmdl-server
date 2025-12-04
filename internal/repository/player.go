package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayerIDByUserID(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*uuid.UUID, error) {
	query := `
		SELECT id
		FROM players
		WHERE user_id = $1
	`

	var id uuid.UUID
	err := pool.QueryRow(ctx, query, userID).Scan(&id)
	if err != nil {
		log.Printf("Unable to fetch resource: %v", err)
		return nil, err
	}
	return &id, nil
}

func CreatePlayer(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*uuid.UUID, error) {
	query := `
		INSERT INTO players (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id
	`

	var id uuid.UUID
	err := pool.QueryRow(ctx, query, userID).Scan(&id)
	if err != nil {
		log.Printf("Unable to insert resource: %v", err)
		return nil, err
	}
	return &id, nil
}
