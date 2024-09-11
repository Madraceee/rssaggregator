package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/madraceee/rssaggregator/internal/database"
)

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Println("Error while decoding params in createFeed - ", err)
		respondWithError(w, 400, "Invalid Input")
		return
	}

	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}

	feed, err := cfg.DB.CreateFeed(context.Background(), args)
	if err != nil {
		log.Println("Error while inserting feed into database", feed, err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJson(w, 201, feed)
}
