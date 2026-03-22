package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func CreateClassicRecord(ctx context.Context, pool *pgxpool.Pool, record model.ClassicRecord) error {
	query := `
		INSERT INTO classic_records (classic_level_id, player_id, record_percentage, completed_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := pool.Exec(ctx, query, record.ClassicLevelID, record.PlayerID, record.RecordPercentage, record.CompletedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("classic record already exists (constraint: %s): %w", pgErr.ConstraintName, err)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("referenced record does not exist (constraint: %s): %w", pgErr.ConstraintName, err)
			case pgerrcode.NotNullViolation:
				return fmt.Errorf("required field missing (column: %s): %w", pgErr.ColumnName, err)
			case pgerrcode.CheckViolation:
				return fmt.Errorf("value failed validation (constraint: %s): %w", pgErr.ConstraintName, err)
			default:
				return fmt.Errorf("database error [%s]: %w", pgErr.Code, err)
			}
		}
		return fmt.Errorf("failed to create classic record: %w", err)
	}
	return nil
}

func GetAllClassicRecordsByClassicLevel(ctx context.Context, pool *pgxpool.Pool, classicLevelId uuid.UUID) ([]*model.ClassicRecord, error) {
	query := `
		SELECT 
            cr.id, cr.player_id, p.id, p.classic_points, r.id, r.classic_points
        FROM classic_records cr
        JOIN players p ON p.id = cr.player_id
        JOIN regions r ON r.id = p.region_id
        WHERE cr.classic_level_id = $1
	`

	rows, err := pool.Query(ctx, query, classicLevelId)
	if err != nil {
		return nil, fmt.Errorf("failed to query classic records: %w", err)
	}
	defer rows.Close()

	var records []*model.ClassicRecord
	for rows.Next() {
		record := &model.ClassicRecord{}
		if err := rows.Scan(&record.ID, &record.PlayerID, &record.Player.ID, &record.Player.ClassicPoints, &record.Player.Region.ID, &record.Player.Region.ClassicPoints); err != nil {
			return nil, fmt.Errorf("failed to scan classic record id: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return records, nil
}
