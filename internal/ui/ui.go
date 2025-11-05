package ui

import (
	"hypermedia/internal/models"
	"net/http"
)

func Handlers(cfg *models.APIConfig) map[string]http.HandlerFunc {
	u := &UI{cfg}

	return map[string]http.HandlerFunc{
		"/":                                     u.RedirectRoot,
		"/favicon.ico":                          u.Favicon,
		"GET /houses":                           u.Houses,
		"GET /houses/{house_id}/flats":          u.HouseView,
		"GET /houses/create":                    u.CreateHouseForm,
		"POST /houses/create":                   u.HandleCreateHouse,
		"GET /houses/{house_id}/flats/create":   u.FlatsCreate,
		"POST /houses/{house_id}/flats":         u.HandleCreateFlats,
		"GET /houses/{house_id}/payment-plans":  u.ServePaymentPlans,
		"POST /houses/{house_id}/payment-plans": u.CreatePaymentPlan,
	}
}
