package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"service/internal/handler"
	"service/internal/repository"
	"service/internal/service"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8080"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var db *sql.DB
	var err error
	// Ожидание поднятия БД в Docker
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			break
		}
		logger.Info("Waiting for database...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Тюнинг пула соединений (защита от перегрузок)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(repo, svc, logger)

	logger.Info("Starting server", "port", serverPort)
	if err := http.ListenAndServe(":"+serverPort, h.InitRoutes()); err != nil {
		logger.Error("Server stopped", "error", err)
	}
}
