package ui

import (
	"hypermedia/internal/models"
	"net/http"
)

func Handlers(cfg *models.APIConfig) map[string]http.HandlerFunc {
	u := &UI{cfg}

	return map[string]http.HandlerFunc{
		"/":                             u.RedirectRoot,
		"GET /houses":                   u.Houses,
		"GET /houses/{home_id}":         u.HouseView,
		"GET /houses/create":            u.CreateHouseForm,
		"POST /houses/create":           u.HandleCreateHouse,
		"POST /houses/{house_id}/flats": u.CreateFlats,
	}
}
