package postgres

import (
	"database/sql"
	"fmt"
	"go-rest-server/internal/config"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgres
func Connect(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	// dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectWithRetry пытается подключиться к Postgres с retry и логами
func ConnectWithRetry(cfg config.Config, logger *slog.Logger) *sql.DB {
	const maxAttempts = 10
	const delay = 2

	var db *sql.DB
	var err error

	for i := 1; i <= maxAttempts; i++ {
		db, err = Connect(cfg)

		if err == nil {
			logger.Info("Подключение к Postgres успешно")
			return db
		}

		logger.Warn("Postgres не готов, повторная попытка...", "attempt", i, "error", err)
		time.Sleep(delay)
	}

	if err != nil {
		logger.Error("Не удалось подключиться к Postgres", "error", err)
		// fmt.Errorf("не удалось подключиться к Postgres после %d попыток: %w", maxAttempts, err)
		os.Exit(1)
	}

	return db

}
