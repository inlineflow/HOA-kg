package ui

import (
	"context"
	"fmt"
	"hypermedia/internal/models"

	// "log"
	"net/http"

	"github.com/a-h/templ"
	// "strconv"
	//
	// "github.com/google/uuid"
)

type UI struct {
	cfg *models.APIConfig
}

func (u *UI) RedirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/houses", http.StatusSeeOther)
}

func (u *UI) Houses(w http.ResponseWriter, r *http.Request) {
	dbHouses, err := u.cfg.DB.GetAllHouses(context.Background())
	if err != nil {
		HandleError(w, r, fmt.Errorf("Error while fetching data from database: %v\n", err), 500)
		return
	}

	houses := models.Map(dbHouses, models.ToHouseVM)
	// houses := make([]models.House, len(dbHouses))
	// for i, v := range dbHouses {
	// 	houses[i] = models.ToHouseVM(v)
	// }

	if r.Header.Get("HX-Request") != "" {
		c := HousesGrid(houses)
		c.Render(context.Background(), w)
		return
	}

	c := ServeHouses(houses)
	c.Render(context.Background(), w)

}

func (u *UI) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(CreateHouse(), opts...).ServeHTTP(w, r)
}

func (u *UI) HomePage(w http.ResponseWriter, r *http.Request) {
	dbHouses, err := u.cfg.DB.GetAllHouses(context.Background())
	if err != nil {
		HandleError(w, r, fmt.Errorf("Error while fetching data from database: %v\n", err), 500)
		return
	}

	houses := make([]models.House, len(dbHouses))
	for i, v := range dbHouses {
		houses[i] = models.ToHouseVM(v)
	}

	ID := r.PathValue("home_id")
	var h models.House
	for _, v := range houses {
		if v.ID.String() == ID {
			h = v
		}
	}

	if (h == models.House{}) {
		w.Write([]byte("Home not found"))
		w.WriteHeader(500)
		return
	}

	if r.Header.Get("HX-Request") != "" {
		c := House(h)
		c.Render(context.Background(), w)
		return
	}

	c := HouseView(h)
	c.Render(context.Background(), w)
}

// func (u *UI) CreateAppartments(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	from, err := strconv.Atoi(r.Form.Get("from"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	to, err := strconv.Atoi(r.Form.Get("to"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
//
// 	appartments := []models.Appartment{}
// 	for i := from; i <= to; i++ {
// 		appartments = append(appartments, models.Appartment{
// 			ID:         uuid.NewString(),
// 			FlatNumber: i,
// 		})
// 	}
//
// 	u.cfg.Appartments = appartments
// 	// fmt.Println(u.cfg.Appartments)
//
// 	http.Redirect(w, r, "/homes/", http.StatusSeeOther)
// }
