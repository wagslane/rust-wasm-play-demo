package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/wagslane/rust-wasm-play-demo/internal/workspacemanager"
)

type config struct {
	WorkspaceManager workspacemanager.WorkspaceManager
}

func (cfg config) startAPI(port string) {
	const dir = "app"

	r := mux.NewRouter()

	r.HandleFunc("/v1/compile", cfg.compileHandler).Methods("POST")
	r.PathPrefix(fmt.Sprintf("/%s/", dir)).Handler(http.StripPrefix(fmt.Sprintf("/%s/", dir), http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler: corsMiddleware(r),
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Printf("Starting server on %v\n", port)
	log.Fatal(srv.ListenAndServe())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
