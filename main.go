package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	cfg := apiConfig{}

	mux := http.NewServeMux()

	mux.Handle("GET /app", http.StripPrefix("/app",
		cfg.middlewareMetricsInc(http.FileServer(http.Dir("./web/")))))

	mux.Handle("GET /app/assets/", http.StripPrefix("/app",
		cfg.middlewareMetricsInc(http.FileServer(http.Dir("./web/")))))

	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /metrics", cfg.metrics)
	mux.HandleFunc("POST /reset", cfg.reset)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Add(cfg.fileServerHits.Load() * -1)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())))
}
