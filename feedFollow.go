package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/madraceee/rssaggregator/internal/database"
)

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&params); err != nil {
		log.Println("Error while decoding in createFeedFollow - ", err)
		respondWithError(w, 400, "Invalid Input")
		return
	}

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), args)
	if err != nil {
		log.Println("Error while create a new feed follow ", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJson(w, 201, feedFollow)
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	id := r.PathValue("feed_follow_id")

	feedFollowID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Invalid Feed Follow id", r.URL.Path)
		respondWithError(w, 400, "Invalid input")
		return
	}

	feedFollow, err := cfg.DB.GetFeedFollowByID(r.Context(), feedFollowID)
	if err != nil {
		log.Println("Error while fetching feed follow", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	if feedFollow.UserID != user.ID {
		log.Println("User unauthorized to delete feed")
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = cfg.DB.DeleteFeedFollowUsingID(r.Context(), feedFollowID)
	if err != nil {
		log.Println("Error while deleting Feed Follow ", err)
		respondWithError(w, 500, "Internal server error")
		return
	}
	respondWithJson(w, 204, struct{}{})
}

func (cfg *apiConfig) getUserFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	data, err := cfg.DB.GetFeedsOfUser(r.Context(), user.ID)
	if err != nil {
		log.Println("Error while fetching Feed Follow of user ", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJson(w, 200, data)
}
