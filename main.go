package main

import (
	"context"
	"fmt"
	"hypermedia/internal/component"
	"hypermedia/internal/models"
	"log"
	"net/http"
)

var dev = true

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func handleContacts(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	contacts := []models.Contact{
		models.Contact{
			ID:    "1",
			First: "John",
			Last:  "Rambo",
			Phone: "1111",
		},
		models.Contact{
			ID:    "2",
			First: "Sylvester",
			Last:  "Stalone",
			Phone: "2222",
		},
	}
	fmt.Println("serving /contacts")
	c := component.Contacts(contacts)
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c.Render(ctx, w)
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

	serveMux.HandleFunc("/contacts", handleContacts)
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// c.Render(context.Background(), os.Stdout)
}
