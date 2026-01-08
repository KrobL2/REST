package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"go-rest-server/internal/config"
	"go-rest-server/internal/repository"
	"go-rest-server/internal/service"
	"go-rest-server/internal/storage/postgres"
	httphandler "go-rest-server/internal/transport/http"
	handler "go-rest-server/internal/transport/http/handlers"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Конфигурация
	cfg := config.Load()

	// Логи
	log := setupLogger(cfg.DBHost)

	// База данных
	db := postgres.ConnectWithRetry(cfg, log)

	defer db.Close()
	log.Info("Подключение к Postgres успешно")

	// --- Применяем миграции Goose ---
	log.Info("Применяем миграции...")
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Error("Ошибка при применении миграций", "error", err)
		os.Exit(1)
	}

	log.Info("Миграции успешно применены")

	// Репозитории и сервисы
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Роутер
	r := httphandler.NewRouter(userHandler)

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
		log.Info("Сервер запущен", "port", 8080)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("ListenAndServe ошибка", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Ошибка при shutdown", "error", err)
	} else {
		log.Info("Сервер остановлен корректно")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		handler := slog.NewTextHandler(os.Stderr, nil)
		log = slog.New(handler)
	case envDev, envProd:
		handler := slog.NewJSONHandler(os.Stdout, nil)
		log = slog.New(handler)
	default:
		handler := slog.NewTextHandler(os.Stderr, nil)
		log = slog.New(handler)
	}

	return log
}
