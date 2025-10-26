package ui

import (
	"hypermedia/internal/models"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

func (u *UI) ServePaymentPlans(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, err, 400)
		return
	}
	dbServicePlans, err := u.cfg.DB.GetPaymentPlansByHouseID(r.Context(), houseID)
	if err != nil {
		HandleError(w, r, err, 400)
		return
	}

	servicePlans := models.Map(dbServicePlans, models.ToPaymentPlanVM)
	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(PaymentPlans(servicePlans, houseID), opts...).ServeHTTP(w, r)

}
