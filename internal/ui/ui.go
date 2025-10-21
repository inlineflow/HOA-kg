package ui

import (
	"hypermedia/internal/models"
	"net/http"
)

func Handlers(cfg *models.APIConfig) map[string]http.HandlerFunc {
	u := &UI{cfg}

	return map[string]http.HandlerFunc{
		"/":                    u.RedirectRoot,
		"GET /homes":           u.Homes,
		"GET /homes/{home_id}": u.HomePage,
	}
}
