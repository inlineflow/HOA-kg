package ui

import (
	"errors"
	"fmt"
	"hypermedia/internal/models"
	"hypermedia/internal/ui/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type UI struct {
	cfg *models.APIConfig
}

func (u *UI) RedirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/houses", http.StatusSeeOther)
}

func (u *UI) Houses(w http.ResponseWriter, r *http.Request) {
	dbHouses, err := u.cfg.DB.GetAllHouses(r.Context())
	if err != nil {
		HandleError(w, r, fmt.Errorf("Error while fetching data from database: %v\n", err), 500)
		return
	}

	houses := models.Map(dbHouses, models.ToHouseVM)

	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(pages.ServeHouses(houses), opts...).ServeHTTP(w, r)
}

func (u *UI) CreateHouseForm(w http.ResponseWriter, r *http.Request) {
	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(pages.CreateHouse(), opts...).ServeHTTP(w, r)
}

func (u *UI) HandleCreateHouse(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to decode form values: %v\n", err), 500)
		return
	}

	address := r.Form.Get("address")

	_, err = u.cfg.DB.InsertHouse(r.Context(), address)
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to create a `House`. Err:%v\n", err), 500)
		return
	}

	http.Redirect(w, r, "/houses", http.StatusSeeOther)
}

func (u *UI) HouseView(w http.ResponseWriter, r *http.Request) {
	dbHouses, err := u.cfg.DB.GetAllHouses(r.Context())
	if err != nil {
		HandleError(w, r, fmt.Errorf("Error while fetching data from database: %v\n", err), 500)
		return
	}

	houses := make([]models.House, len(dbHouses))
	for i, v := range dbHouses {
		houses[i] = models.ToHouseVM(v)
	}

	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, &models.PathValueParseError{ResourceKey: "house_id", ParseError: err}, 500)
		return
	}
	var h models.House
	for _, v := range houses {
		if v.ID == houseID {
			h = v
		}
	}

	if (h == models.House{}) {
		HandleError(w, r, errors.New("House not found"), 500)
		return
	}

	dbFlats, err := u.cfg.DB.GetFlatsForHouse(r.Context(), houseID)
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to get `[]Flat` for [house_id:%v]: %v", houseID, err), 500)
		return
	}

	flats := models.Map(dbFlats, models.ToFlatVM)

	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(pages.HouseView(h, flats), opts...).ServeHTTP(w, r)

}
