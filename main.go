package main

import (
	// "context"
	"context"
	"fmt"
	"hypermedia/internal/component"
	"log"
	"net/http"
	// "os"
	// "github.com/a-h/templ"
)

var dev = true

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving root")
	c := component.Hello("Big")
	c.Render(context.Background(), w)
}

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
	// c := component.Hello("John")
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handleRoot)
	serveMux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))))
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// c.Render(context.Background(), os.Stdout)
}
