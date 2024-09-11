package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func (cfg *apiConfig) getUsersByApiKey(w http.ResponseWriter, r *http.Request) {
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

	respondWithJson(w, 200, user)
}
