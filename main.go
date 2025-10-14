package main

import (
	"context"
	"fmt"
	"hypermedia/internal/component"
	"hypermedia/internal/models"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var dev = true
var contacts []models.Contact

// contacts := []models.Contact{
// 	{
// 		ID:    "1",
// 		First: "John",
// 		Last:  "Rambo",
// 		Phone: "1111",
// 	},
// 	{
// 		ID:    "2",
// 		First: "Sylvester",
// 		Last:  "Stalone",
// 		Phone: "2222",
// 	},
// }

func handlePartialContacts(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c := component.ContactsFormList(contacts)
	c.Render(ctx, w)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func handleGetContacts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(contacts)
	searchTerm := r.URL.Query().Get("q")
	fmt.Println("serving /contacts")
	c := component.GetContacts(contacts)
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c.Render(ctx, w)
}

func handleGetContactsNew(w http.ResponseWriter, r *http.Request) {
	c := component.NewContact(models.Contact{})
	c.Render(context.Background(), w)
}

func handlePostContactsNew(w http.ResponseWriter, r *http.Request) {
	m := models.Contact{
		First: r.FormValue("first_name"),
		Last:  r.FormValue("last_name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
		ID:    uuid.NewString(),
	}

	contacts = append(contacts, m)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)

}

func handleGetContactByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("contact_id")
	var c models.Contact
	for _, v := range contacts {
		if v.ID == id {
			c = v
			break
		}
	}

	if c.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Error, contact not found"))
		return
	}

	view := component.ContactDetails(c)
	view.Render(context.Background(), w)

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

func handleContactEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("contact_id")
	var c models.Contact
	for _, v := range contacts {
		if v.ID == id {
			c = v
			break
		}
	}

	if c.ID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Error, contact not found"))
		return
	}

	view := component.EditContact(c)
	view.Render(context.Background(), w)
}

func main() {
	// c := component.Hello("John")
	contacts = []models.Contact{
		{
			ID:    uuid.NewString(),
			First: "John",
			Last:  "Rambo",
			Phone: "1111",
		},
		{
			ID:    uuid.NewString(),
			First: "Sylvester",
			Last:  "Stalone",
			Phone: "2222",
		},
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handleRoot)
	serveMux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))))

	serveMux.HandleFunc("GET /contacts", handleGetContacts)
	serveMux.HandleFunc("GET /contacts/new", handleGetContactsNew)
	serveMux.HandleFunc("POST /contacts/new", handlePostContactsNew)
	serveMux.HandleFunc("GET /contacts/{contact_id}", handleGetContactByID)
	serveMux.HandleFunc("GET /contacts/{contact_id}/edit", handleContactEdit)
	serveMux.HandleFunc("GET /partials/contacts", handlePartialContacts)
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// c.Render(context.Background(), os.Stdout)
}
