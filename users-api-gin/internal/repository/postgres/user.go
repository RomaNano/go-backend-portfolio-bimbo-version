package postgres

import (
	"database/sql"

	"users-api-gin/internal/model"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id, created_at
	`

	return r.db.
		QueryRow(query, user.Email).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *PostgresUserRepository) GetByID(id int64) (*model.User, error) {
	query := `
		SELECT id, email, created_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.
		QueryRow(query, id).
		Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) List() ([]model.User, error) {
	query := `
		SELECT id, email, created_at
		FROM users
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
