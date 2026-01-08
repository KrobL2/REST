package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"go-rest-server/internal/config"
	"go-rest-server/internal/repository"
	"go-rest-server/internal/service"
	"go-rest-server/internal/transport/handler"
)

func main() {
	// --- Логи ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// --- Конфигурация ---
	cfg := config.Load()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// --- Подключение к БД с retry ---
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			break
		}

		slog.Warn("Postgres не готов, повторная попытка...", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		slog.Error("Не удалось подключиться к Postgres", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	slog.Info("Подключение к Postgres успешно")

	// --- Применяем миграции Goose ---
	slog.Info("Применяем миграции...")
	if err := goose.Up(db, "./migrations"); err != nil {
		slog.Error("Ошибка при применении миграций", "error", err)
		os.Exit(1)
	}

	slog.Info("Миграции успешно применены")

	// --- Репозитории и хендлеры ---
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// --- HTTP роутер ---
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", userHandler.GetUsers)
	})

	// --- HTTP сервер с таймаутами ---
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// --- Graceful shutdown ---
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Сервер запущен", "port", 8080)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe ошибка", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Ошибка при shutdown", "error", err)
	} else {
		slog.Info("Сервер остановлен корректно")
	}
}
