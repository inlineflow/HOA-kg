package ui

import (
	"context"
	"hypermedia/internal/models"
	"net/http"
)

type UI struct {
	cfg *models.APIConfig
}

func (u *UI) ServeRoot(w http.ResponseWriter, r *http.Request) {
	c := Root()
	c.Render(context.Background(), w)
}
