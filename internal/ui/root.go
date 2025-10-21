package ui

import (
	"context"
	"hypermedia/internal/models"
	"net/http"
)

type UI struct {
	cfg *models.APIConfig
}

func (u *UI) RedirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/homes", http.StatusSeeOther)
	// c := Root()
	// c.Render(context.Background(), w)
}

func (u *UI) Homes(w http.ResponseWriter, r *http.Request) {
	c := ServeHomes(u.cfg.Data)
	c.Render(context.Background(), w)
}

func (u *UI) HomePage(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("home_id")
	var h models.Home
	for _, v := range u.cfg.Data {
		if v.ID == ID {
			h = v
		}
	}

	if (h == models.Home{}) {
		w.Write([]byte("Home not found"))
		w.WriteHeader(500)
		return
	}
	c := HomeView(h)
	c.Render(context.Background(), w)
}
