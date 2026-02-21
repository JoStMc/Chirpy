package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
} 

func main() {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /reset", cfg.handlerReset)

	server := http.Server{
	    Handler: mux,
		Addr: ":8080",
	} 
	log.Println("Serving on port 8080")

	log.Fatal(server.ListenAndServe())
}
