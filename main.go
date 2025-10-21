package main

import (
	"fmt"
	"hypermedia/internal/models"
	"hypermedia/internal/ui"
	"log"
	"net/http"
)

var dev = true

func disableCacheInDevMode(next http.Handler) http.Handler {
	if !dev {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func main() {

	cfg := &models.APIConfig{}
	mux := http.NewServeMux()
	mux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))))

	for pattern, handler := range ui.Handlers(cfg) {
		mux.HandleFunc(pattern, handler)
	}

	server := http.Server{Handler: mux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
