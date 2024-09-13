package main

import (
	"log"
	"net/http"

	"github.com/madraceee/rssaggregator/internal/database"
)

func (cfg *apiConfig) getPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Println("Error while fetching posts for user", user.ID, err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJson(w, 200, posts)
}
