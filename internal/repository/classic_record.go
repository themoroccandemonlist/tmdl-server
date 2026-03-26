package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/dto"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func AdminGetAllClassicRecords(ctx context.Context, pool *pgxpool.Pool, limit int, offset int, search string) ([]dto.AdminClassicRecordRow, error) {
	query := `
		SELECT r.id, p.id, p.username, cl.id, cl.name, r.progress, r.completed_at, r.device, r.footage, r.raw_footage
		FROM classic_records r
		JOIN players p ON p.id = r.player_id
		JOIN classic_levels cl ON cl.id = r.classic_level_id
		WHERE (p.username ILIKE '%' || $3 || '%' OR cl.name ILIKE '%' || $3 || '%')
		ORDER BY r.completed_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := pool.Query(ctx, query, limit, offset, search)
	if err != nil {
		return nil, fmt.Errorf("failed to query classic records: %w", err)
	}
	defer rows.Close()

	var records []dto.AdminClassicRecordRow
	for rows.Next() {
		var row dto.AdminClassicRecordRow
		if err := rows.Scan(&row.ID, &row.Player.ID, &row.Player.Username, &row.ClassicLevel.ID, &row.ClassicLevel.Name, &row.Progress, &row.Date, &row.Device, &row.Footage, &row.RawFootage); err != nil {
			return nil, fmt.Errorf("failed to scan classic record row: %w", err)
		}
		records = append(records, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return records, nil
}

func GetClassicRecordsCount(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM classic_records
	`

	var count int
	if err := pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count classic records: %w", err)
	}
	return count, nil
}

func CreateClassicRecord(ctx context.Context, pool *pgxpool.Pool, record model.ClassicRecord) error {
	query := `
		INSERT INTO classic_records (classic_level_id, player_id, record_percentage, completed_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := pool.Exec(ctx, query, record.ClassicLevelID, record.PlayerID, record.Progress, record.CompletedAt)
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
