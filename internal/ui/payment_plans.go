package ui

import (
	"fmt"
	"hypermedia/internal/database"
	"hypermedia/internal/models"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5/pgtype"
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

	dbPeriods, err := u.cfg.DB.GetPeriods(r.Context())
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to get dbPeriods. Error: %v", err), 400)
		return
	}
	periods := models.Map(dbPeriods, models.ToPeriodVM)

	props := PaymentPlanFormProps{
		houseID: houseID,
		periods: periods,
	}
	fmt.Println(props)
	fmt.Println(servicePlans)
	// w.Write([]byte("hello world"))
	templ.Handler(PaymentPlans(servicePlans, props), opts...).ServeHTTP(w, r)

}

func (u *UI) CreatePaymentPlan(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, err, 400)
		return
	}

	err = r.ParseForm()
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to decode form values: %v\n", err), 500)
		return
	}
	periodID, err := uuid.Parse(r.FormValue("period_id"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to parse periodID: %v\n", err), 500)
		return
	}
	amount, err := decimal.Parse(r.FormValue("amount"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to parse amount: %v\n", err), 500)
		return
	}

	pp, err := u.cfg.DB.CreatePaymentPlan(r.Context(), database.CreatePaymentPlanParams{
		HouseID:  houseID,
		PeriodID: periodID,
		DateFrom: database.Now(),
		DateTo: pgtype.Timestamp{
			Valid: false,
		},
		MonthlyAmount: amount,
	})
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to create a payment plan: %v\n", err), 500)
		return
	}

	fmt.Println(pp)
	http.Redirect(w, r, fmt.Sprintf("/houses/%s/payment-plans", houseID), http.StatusSeeOther)
	// w.Write([]byte("hello world"))
}
