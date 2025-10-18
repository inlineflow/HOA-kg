package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hypermedia/internal/component"
	"hypermedia/internal/models"
	"hypermedia/internal/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var dev = true

func filterContacts(q string, c []models.Contact) []models.Contact {
	result := []models.Contact{}
	for _, v := range c {
		if v.Email == q || v.First == q || v.Last == q || v.Phone == q {
			result = append(result, v)
		}
	}

	return result
}

// func (cfg *APIConfig) handlePartialContacts(w http.ResponseWriter, r *http.Request) {
func (cfg *APIConfig) handlePartialContacts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	contacts := cfg.db.GetContacts()
	i := (page - 1) * 10
	j := min(i+10, len(contacts))
	c := component.ContactList(contacts[i:j], page)
	c.Render(context.Background(), w)
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
	var contacts []models.Contact
	if searchTerm != "" {
		contacts = filterContacts(searchTerm, cfg.db.GetContacts())
	} else {
		contacts = cfg.db.GetContacts()
	}
	i := (page - 1) * 10
	j := min(i+10, len(contacts))
	if r.Header.Get("HX-Trigger") == "search" {
		c := component.ContactList(contacts, page)
		c.Render(context.Background(), w)
		return
	}

	c := component.GetContacts(contacts[i:j], page, cfg.archiver)
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

func (db *DB) BulkDelete(ids []string) {
	idsToDelete := make(map[string]struct{})
	for _, id := range ids {
		idsToDelete[id] = struct{}{}
	}

	newData := []models.Contact{}
	for _, v := range db.Data {
		if _, exists := idsToDelete[v.ID]; !exists {
			newData = append(newData, v)
		}
	}

	db.Data = newData
	writeContacts(db.Data)
}

type DBX interface {
	GetContacts() []models.Contact
	AddContact(c models.Contact)
	UpdateContact(c models.Contact)
	DeleteContact(id string)
	BulkDelete(ids []string)
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
	db       DBX
	archiver *services.Archiver
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

func (cfg *APIConfig) handleDeleteContactByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling delete")
	id := r.PathValue("contact_id")
	cfg.db.DeleteContact(id)
	if r.Header.Get("HX-Trigger") == "delete-form" {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	}

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

func (cfg *APIConfig) handleContactCount(w http.ResponseWriter, r *http.Request) {
	countStr := fmt.Sprintf("%d total contacts", len(cfg.db.GetContacts()))
	w.Write([]byte(countStr))
}

func (cfg *APIConfig) handleDeleteContacts(w http.ResponseWriter, r *http.Request) {
	selectedIDsStr := r.URL.Query().Get("selected_contact_ids")

	searchTerm := r.URL.Query().Get("q")
	IDs := strings.Split(selectedIDsStr, ",")
	fmt.Println(IDs)
	// for _, v := range cfg.db.GetContacts() {
	// 	for _, j := range IDs {
	// 		if v.ID == j {
	// 			fmt.Println("Matched on ", j)
	// 			cfg.db.DeleteContact(j)
	// 		}
	// 	}
	// }
	cfg.db.BulkDelete(IDs)
	ctx := context.WithValue(context.Background(), "search_term", searchTerm)
	c := component.GetContacts(cfg.db.GetContacts(), 1, cfg.archiver)
	c.Render(ctx, w)
	// fmt.Println(r)
	// fmt.Println(r.Form)
	// fmt.Println("query: ", r.URL.Query())
	// fmt.Println("r.PostForm: ", r.PostForm)
	// fmt.Println(IDs)
}

func (cfg *APIConfig) handlePostContactsArchive(w http.ResponseWriter, r *http.Request) {
	cfg.archiver.Start()
	time.Sleep(100 * time.Millisecond)
	c := component.ArchiveDownloadButton(cfg.archiver)
	c.Render(context.Background(), w)
}

func (cfg *APIConfig) handleGetContactsArchive(w http.ResponseWriter, r *http.Request) {
	c := component.ArchiveDownloadButton(cfg.archiver)
	c.Render(context.Background(), w)
}

func main() {
	// c := component.Hello("John")
	contacts, err := loadContacts()
	if err != nil {
		log.Fatal(err)
	}

	db := DB{contacts}
	cfg := &APIConfig{&db, services.NewArchiver()}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handleRoot)
	serveMux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets", http.FileServer(http.Dir("assets")))))

	serveMux.HandleFunc("GET /contacts", cfg.handleGetContacts)
	serveMux.HandleFunc("DELETE /contacts", cfg.handleDeleteContacts)
	serveMux.HandleFunc("GET /contacts/new", cfg.handleGetContactsNew)
	serveMux.HandleFunc("POST /contacts/new", cfg.handlePostContactsNew)
	serveMux.HandleFunc("GET /contacts/{contact_id}/email", cfg.handleValidateEmail)
	serveMux.HandleFunc("GET /contacts/{contact_id}", cfg.handleGetContactByID)
	serveMux.HandleFunc("GET /contacts/{contact_id}/edit", cfg.handleContactEdit)
	serveMux.HandleFunc("POST /contacts/{contact_id}/edit", cfg.handlePostEditContact)
	serveMux.HandleFunc("DELETE /contacts/{contact_id}", cfg.handleDeleteContactByID)
	serveMux.HandleFunc("GET /partials/contacts", cfg.handlePartialContacts)
	serveMux.HandleFunc("GET /contacts/count", cfg.handleContactCount)
	serveMux.HandleFunc("POST /contacts/archive", cfg.handlePostContactsArchive)
	serveMux.HandleFunc("GET /contacts/archive", cfg.handleGetContactsArchive)
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	fmt.Println("Started on localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// c.Render(context.Background(), os.Stdout)
}
