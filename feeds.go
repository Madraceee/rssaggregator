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

	feedArgs := database.CreateFeedParams{

		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}

	feed, err := cfg.DB.CreateFeed(context.Background(), feedArgs)
	if err != nil {
		log.Println("Error while inserting feed into database", feed, err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	// Creating Feed Follow for the user
	feedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowArgs)
	if err != nil {
		log.Println("Error while create a new feed follow ", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	result := struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	}

	respondWithJson(w, 201, result)
}

func (cfg *apiConfig) getAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(context.Background())
	if err != nil {
		log.Println("Error while fetching feeds", err)
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	respondWithJson(w, 200, feeds)
}
