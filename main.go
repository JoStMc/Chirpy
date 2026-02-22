package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/JoStMc/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
} 

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load env:", err)
	} 
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("unable to open db:", err)
	} 

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
	    Handler: mux,
		Addr: ":8080",
	} 
	log.Println("Serving on port 8080")

	log.Fatal(server.ListenAndServe())
}
