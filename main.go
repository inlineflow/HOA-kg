package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hypermedia/internal/component"
	"hypermedia/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var dev = true

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

// func (cfg *APIConfig) handlePartialContacts(w http.ResponseWriter, r *http.Request) {
func (cfg *APIConfig) handlePartialContacts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	contacts := cfg.db.GetContacts()
	i := (page - 1) * 10
	j := min(i+10, len(contacts))
	// if j > len(contacts) {
	// 	j = len(contacts)
	// }
	searchTerm := r.URL.Query().Get("q")
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c := component.ContactList(contacts[i:j], page)
	c.Render(ctx, w)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (cfg *APIConfig) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	searchTerm := r.URL.Query().Get("q")
	i := (page - 1) * 10
	j := i + 10
	c := component.GetContacts(cfg.db.GetContacts()[i:j], page)
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c.Render(ctx, w)
}

func (cfg *APIConfig) handleGetContactsNew(w http.ResponseWriter, r *http.Request) {
	c := component.NewContact(models.Contact{})
	c.Render(context.Background(), w)
}

func (cfg *APIConfig) handlePostContactsNew(w http.ResponseWriter, r *http.Request) {
	m := models.Contact{
		First: r.FormValue("first_name"),
		Last:  r.FormValue("last_name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
		ID:    uuid.NewString(),
	}

	cfg.db.AddContact(m)
	// contacts = append(contacts, m)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)

}

func (cfg *APIConfig) handleGetContactByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("contact_id")
	var c models.Contact
	for _, v := range cfg.db.GetContacts() {
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

func (cfg *APIConfig) handleContactEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("contact_id")
	var c models.Contact
	for _, v := range cfg.db.GetContacts() {
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

type DB struct {
	Data []models.Contact
}

func (db *DB) GetContacts() []models.Contact {
	return db.Data
}

func (db *DB) UpdateContact(newC models.Contact) {
	for i, v := range db.Data {
		if v.ID == newC.ID {
			db.Data[i] = newC
		}
	}

	writeContacts(db.Data)

}

func (db *DB) DeleteContact(id string) {
	i := -1
	for j, v := range db.Data {
		if v.ID == id {
			i = j
		}
	}

	db.Data = append(db.Data[:i], db.Data[i+1:]...)
	writeContacts(db.Data)
}

func (db *DB) AddContact(c models.Contact) {
	db.Data = append(db.Data, c)
	writeContacts(db.Data)
}

func (db *DB) ReloadContactDB(contacts []models.Contact) {
	db.Data = contacts
}

type DBX interface {
	GetContacts() []models.Contact
	AddContact(c models.Contact)
	UpdateContact(c models.Contact)
	DeleteContact(id string)
	ReloadContactDB(contacts []models.Contact)
}

func loadContacts() ([]models.Contact, error) {
	fileBytes, err := os.ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	var contacts []models.Contact
	err = json.Unmarshal(fileBytes, &contacts)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func writeContacts(contacts []models.Contact) error {
	f, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(contacts)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type APIConfig struct {
	db DBX
}

func (cfg *APIConfig) writeAndLoadContacts(contacts []models.Contact) error {
	err := writeContacts(contacts)
	if err != nil {
		return fmt.Errorf("Error while writing contacts: %v", err)
	}

	c, err := loadContacts()
	if err != nil {
		return fmt.Errorf("Error while loading contacts after write: %v", err)
	}

	cfg.db.ReloadContactDB(c)
	return nil
}

func (cfg *APIConfig) handlePostEditContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("contact_id")
	m := models.Contact{
		First: r.FormValue("first_name"),
		Last:  r.FormValue("last_name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
		ID:    id,
	}

	cfg.db.UpdateContact(m)
	http.Redirect(w, r, fmt.Sprintf("/contacts/%s", m.ID), http.StatusSeeOther)
}

func (cfg *APIConfig) handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling delete")
	id := r.PathValue("contact_id")
	cfg.db.DeleteContact(id)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (cfg *APIConfig) handleValidateEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("validating email")
	newEmail := r.URL.Query().Get("email")
	if newEmail == "" {
		return
	}
	fmt.Println("new email: ", newEmail)
	errs := []string{}
	for _, v := range cfg.db.GetContacts() {
		if v.Email == newEmail {
			errs = append(errs, "This email is already registered.")
		}
	}

	if len(errs) > 0 {
		w.Write([]byte(strings.Join(errs, " ")))
	}
}

func main() {
	// c := component.Hello("John")
	contacts, err := loadContacts()
	if err != nil {
		log.Fatal(err)
	}

	db := DB{contacts}
	cfg := &APIConfig{&db}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handleRoot)
	serveMux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))))

	serveMux.HandleFunc("GET /contacts", cfg.handleGetContacts)
	serveMux.HandleFunc("GET /contacts/new", cfg.handleGetContactsNew)
	serveMux.HandleFunc("POST /contacts/new", cfg.handlePostContactsNew)
	serveMux.HandleFunc("GET /contacts/{contact_id}/email", cfg.handleValidateEmail)
	serveMux.HandleFunc("GET /contacts/{contact_id}", cfg.handleGetContactByID)
	serveMux.HandleFunc("GET /contacts/{contact_id}/edit", cfg.handleContactEdit)
	serveMux.HandleFunc("POST /contacts/{contact_id}/edit", cfg.handlePostEditContact)
	serveMux.HandleFunc("DELETE /contacts/{contact_id}", cfg.handleDeleteContact)
	serveMux.HandleFunc("GET /partials/contacts", cfg.handlePartialContacts)
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// c.Render(context.Background(), os.Stdout)
}
