package main

import (
	"database/sql"

	"github.com/madraceee/rssaggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func newApiConfig(db *sql.DB) *apiConfig {
	return &apiConfig{
		DB: database.New(db),
	}
}
