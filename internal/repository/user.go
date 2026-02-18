package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themoroccandemonlist/tmdl-server/internal/model"
)

func GetUserByEmailAndSub(ctx context.Context, pool *pgxpool.Pool, email string, sub string) (*model.User, error) {
	query := `
		SELECT u.id, u.email, u.sub, ARRAY_AGG(r.name) AS role_names, u.is_banned, u.is_deleted
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE u.email = $1 AND u.sub = $2
		GROUP BY u.id, u.email, u.sub, u.is_banned, u.is_deleted;
	`

	var user model.User
	err := pool.QueryRow(ctx, query, email, sub).Scan(&user.ID, &user.Email, &user.Sub, &user.Roles, &user.IsBanned, &user.IsDeleted)
	if err != nil {
		log.Printf("Unable to fetch resource: %v", err)
		return nil, err
	}
	return &user, nil
}

func CreateUser(ctx context.Context, pool *pgxpool.Pool, email string, sub string) (*model.User, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var user model.User

	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, sub)
		VALUES ($1, $2)
		RETURNING id, email, sub, is_banned, is_deleted
	`, email, sub).Scan(&user.ID, &user.Email, &user.Sub, &user.IsBanned, &user.IsDeleted)
	if err != nil {
		log.Printf("Unable to insert user: %v", err)
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = 'USER'
	`, user.ID)
	if err != nil {
		log.Printf("Unable to assign role: %v", err)
		return nil, err
	}

	err = tx.QueryRow(ctx, `
		SELECT u.id, u.email, u.sub, ARRAY_REMOVE(ARRAY_AGG(r.name), NULL) AS role_names, u.is_banned, u.is_deleted
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE u.id = $1
		GROUP BY u.id, u.email, u.sub, u.is_banned, u.is_deleted
	`, user.ID).Scan(&user.ID, &user.Email, &user.Sub, &user.Roles, &user.IsBanned, &user.IsDeleted)
	if err != nil {
		log.Printf("Unable to fetch user: %v", err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	log.Println(user)
	return &user, nil
}