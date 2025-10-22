package models

import (
	"context"
	"errors"
	"fmt"
	"hypermedia/internal/database"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type APIConfig struct {
	Logger *zap.SugaredLogger
	DB     *database.Queries
	Pool   *pgxpool.Pool
}

func NewConfig() (*APIConfig, error) {
	logger, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL missing in env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbURL)
	fmt.Println(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v\n", err)
	}
	db := database.New(pool)

	return &APIConfig{Logger: sugar, DB: db, Pool: pool}, nil
}
