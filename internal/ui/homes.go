package ui

import (
	"context"
	"hypermedia/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type UI struct {
	cfg *models.APIConfig
}

func (u *UI) RedirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/homes", http.StatusSeeOther)
}

func (u *UI) Homes(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		c := HomesGrid(u.cfg.Homes)
		c.Render(context.Background(), w)
		return
	}

	c := ServeHomes(u.cfg.Homes)
	c.Render(context.Background(), w)

}

func (u *UI) HomePage(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("home_id")
	var h models.Home
	for _, v := range u.cfg.Homes {
		if v.ID == ID {
			h = v
		}
	}

	if (h == models.Home{}) {
		w.Write([]byte("Home not found"))
		w.WriteHeader(500)
		return
	}

	if r.Header.Get("HX-Request") != "" {
		c := Home(h)
		c.Render(context.Background(), w)
		return
	}

	c := HomeView(h)
	c.Render(context.Background(), w)
}

func (u *UI) CreateAppartments(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil {
		log.Println(err)
		return
	}
	to, err := strconv.Atoi(r.Form.Get("to"))
	if err != nil {
		log.Println(err)
		return
	}

	apppartments := []models.Appartment{}
	for i := from; i <= to; i++ {
		apppartments = append(apppartments, models.Appartment{
			ID:         uuid.NewString(),
			FlatNumber: i,
		})
	}

	u.cfg.Appartments = apppartments
	// fmt.Println(u.cfg.Appartments)

	http.Redirect(w, r, "/homes/", http.StatusSeeOther)
}
