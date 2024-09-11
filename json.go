package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	// Marshlling payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln("Error while marshling json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payloadBytes)
	return
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	errorPayload := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	respondWithJson(w, code, errorPayload)
}
