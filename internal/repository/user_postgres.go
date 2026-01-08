package repository

import (
	"database/sql"
	"go-rest-server/internal/domain"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *PostgresUserRepo) Create(user domain.User) (int, error) {
	var id int

	err := r.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
