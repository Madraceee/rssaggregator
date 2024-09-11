package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func health(w http.ResponseWriter, r *http.Request) {
	result := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJson(w, 200, result)
}

func errFunc(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func main() {
	godotenv.Load("./.env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	if port == "" || dbURL == "" {
		log.Fatalln("Env variables not present")
	}

	log.Println("Connecting to database")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Error connecting to Database")
	}
	defer db.Close()
	log.Println("Connected to database")

	apiConfig := newApiConfig(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", health)
	mux.HandleFunc("GET /v1/err", errFunc)

	mux.HandleFunc("POST /v1/users", apiConfig.addUser)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.getUser))

	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.createFeed))

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Println("Server starting at :", port)
	server.ListenAndServe()
}
