package postgres

import (
	"database/sql"

	"users-api-http/internal/model"
)

type UserRepository struct{
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository{
	return &UserRepository{db: db}
}


func (r *UserRepository) Create(user *model.User) error {
	query := `
	INSERT INTO users (email)
	VALUE ($1)
	RETURNING id, created_at
	`
	return r.db.
	QueryRow(query, user.Email).
	Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
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
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *UserRepository) List() ([]model.User, error) {
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
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}


