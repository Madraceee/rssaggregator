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

type Feed struct {
	Id            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserId        uuid.UUID `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at,omitempty"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	newFeed := Feed{
		Id:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserId:    feed.UserID,
	}

	if feed.LastFetchedAt.Valid {
		newFeed.LastFetchedAt = feed.LastFetchedAt.Time
	}

	return newFeed
}

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
		Feed       Feed                `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       databaseFeedToFeed(feed),
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

	updatedFeeds := make([]Feed, len(feeds))
	for k, v := range feeds {
		updatedFeeds[k] = databaseFeedToFeed(v)
	}

	respondWithJson(w, 200, updatedFeeds)
}
