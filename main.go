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
	jwtSecret string
	polkaKey string
	dbQueries *database.Queries
	fileserverHits atomic.Int32
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
		jwtSecret: os.Getenv("TOKEN_SECRET"),
		polkaKey: os.Getenv("POLKA_KEY"),
		dbQueries: database.New(db),
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	cfg.registerHandlers(mux)


	server := http.Server{
	    Handler: mux,
		Addr: ":8080",
	} 
	log.Println("Serving on port 8080")

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) registerHandlers(mux *http.ServeMux) {
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeUser)

	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
}
