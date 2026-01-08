package repository

import (
	"database/sql"
	"fmt"
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

	fmt.Print(users)

	return users, nil
}

func (r *PostgresUserRepo) Create(user domain.User) (domain.User, error) {
	err := r.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
