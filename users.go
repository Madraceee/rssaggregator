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

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Name string `json:"name"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	body := requestBody{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Println(err)
		respondWithError(w, 400, "Invalid Input")
		return
	}

	apiKey, err := randomHex(32)
	if err != nil {
		log.Println("Error while generating apiKey - ", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	log.Printf("User %s has a apiKey of %s\n", body.Name, apiKey)
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
		ApiKey:    apiKey,
	}

	result, err := cfg.DB.CreateUser(context.Background(), args)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Server Error")
		return
	}

	respondWithJson(w, 201, result)
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, user)
}
