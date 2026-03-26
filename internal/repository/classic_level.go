package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/dto"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func AdminGetAllClassicLevels(ctx context.Context, pool *pgxpool.Pool, limit int, offset int, search string) ([]dto.AdminClassicLevelRow, error) {
	query := `
		SELECT id, name, publisher, difficulty, ranking, points
		FROM classic_levels
		WHERE (name ILIKE '%' || $3 || '%' OR publisher ILIKE '%' || $3 || '%')
		ORDER BY ranking
		LIMIT $1 OFFSET $2
	`

	rows, err := pool.Query(ctx, query, limit, offset, search)
	if err != nil {
		return nil, fmt.Errorf("failed to query classic levels: %w", err)
	}
	defer rows.Close()

	var levels []dto.AdminClassicLevelRow
	for rows.Next() {
		var row dto.AdminClassicLevelRow
		if err := rows.Scan(&row.ID, &row.Name, &row.Publisher, &row.Difficulty, &row.Ranking, &row.Points); err != nil {
			return nil, fmt.Errorf("failed to scan classic level row: %w", err)
		}
		levels = append(levels, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return levels, nil
}

func GetClassicLevelsCount(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM classic_levels
	`

	var count int
	if err := pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count classic levels: %w", err)
	}
	return count, nil
}

func CreateClassicLevel(ctx context.Context, pool *pgxpool.Pool, level model.ClassicLevel) error {
	query := `
		INSERT INTO classic_levels (level_id, name, publisher, difficulty, duration, ranking, list_percentage, points, minimum_points, youtube_link, thumbnail_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := pool.Exec(ctx, query, level.LevelID, level.Name, level.Publisher, level.Difficulty, level.Duration, level.Ranking, level.ListPercentage, level.Points, level.MinimumPoints, level.YoutubeLink, level.ThumbnailPath)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("classic level already exists (constraint: %s): %w", pgErr.ConstraintName, err)
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
		return fmt.Errorf("failed to create classic level: %w", err)
	}
	return nil
}

func GetAffectedClassicLevels(ctx context.Context, pool *pgxpool.Pool, ranking int) ([]*model.ClassicLevel, error) {
	query := `
		SELECT id, ranking
		FROM classic_levels
		WHERE ranking >= $1
	`

	rows, err := pool.Query(ctx, query, ranking)
	if err != nil {
		return nil, fmt.Errorf("failed to query affected classic levels: %w", err)
	}
	defer rows.Close()

	var levels []*model.ClassicLevel
	for rows.Next() {
		level := &model.ClassicLevel{}
		if err := rows.Scan(&level.ID, &level.Ranking); err != nil {
			return nil, fmt.Errorf("failed to scan classic level id: %w", err)
		}
		levels = append(levels, level)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return levels, nil
}

// to verify
func UpdateClassicLevelsPoints(ctx context.Context, pool *pgxpool.Pool, levels []*model.ClassicLevel) error {
    ids := make([]uuid.UUID, len(levels))
    points := make([]decimal.Decimal, len(levels))
    minimumPoints := make([]decimal.Decimal, len(levels))
    rankings := make([]int, len(levels))

    for i, l := range levels {
        ids[i] = l.ID
        points[i] = l.Points
        minimumPoints[i] = l.MinimumPoints
        rankings[i] = l.Ranking
    }

    _, err := pool.Exec(ctx, `
        UPDATE classic_levels SET
            points = data.points,
            minimum_points = data.minimum_points,
            ranking = data.ranking
        FROM unnest($1::uuid[], $2::numeric[], $3::numeric[], $4::int[]) AS data(id, points, minimum_points, ranking)
        WHERE classic_levels.id = data.id
    `, ids, points, minimumPoints, rankings)
    if err != nil {
        return fmt.Errorf("failed to bulk update classic levels: %w", err)
    }
    return nil
}