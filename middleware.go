package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/madraceee/rssaggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			log.Println("Empty Header")
			respondWithError(w, 400, "Invalid Inputs")
			return
		}

		authHeader := strings.Split(authorization, " ")
		if len(authHeader) != 2 {
			log.Println("Invalid Header - ", authorization)
			respondWithError(w, 400, "Invalid Inputs")
			return
		}

		if authHeader[0] != "ApiKey" {
			log.Println("Invalid Header - ", authorization)
			respondWithError(w, 400, "Invalid Inputs")
			return
		}

		apiKey := authHeader[1]
		user, err := cfg.DB.FetchByApiKey(context.Background(), apiKey)
		if err != nil {
			log.Println("Error while fetching user using ApiKey - ", err)
			respondWithError(w, 500, "Internal Server Error")
			return
		}
		handler(w, r, user)
	}
}
