package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db DBQuerier
}

const usersTable = "users"

func CreateUserRepository(db DBQuerier) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, userModel *user.User) (*user.User, error) {
	const query = `
		INSERT INTO %s (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, created_at
	`
	sql := fmt.Sprintf(query, usersTable)

	newUser := &user.User{}
	err := r.db.QueryRow(ctx, sql, userModel.ID, userModel.Email, userModel.PasswordHash).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, fmt.Errorf("failed to create userModel: %w", err)
	}
	return newUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*user.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
        FROM %s
		WHERE id = $1
	`
	sql := fmt.Sprintf(query, usersTable)

	userModel := &user.User{}
	err := r.db.QueryRow(ctx, sql, userId).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get userModel: %w", err)
	}
	return userModel, nil
}

func (r *UserRepository) GetUserByProfileId(ctx context.Context, profileId uuid.UUID) (*user.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
        FROM %s
		WHERE id = (SELECT user_id FROM %s WHERE id = $1)
	`
	sql := fmt.Sprintf(query, usersTable, profilesTable)
	user := &user.User{}
	err := r.db.QueryRow(ctx, sql, profileId).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
