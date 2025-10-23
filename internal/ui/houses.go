package ui

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"hypermedia/internal/database"
	"hypermedia/internal/models"
	"net/http"
	"strconv"
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

	templ.Handler(ServeHouses(houses), opts...).ServeHTTP(w, r)
}

func (u *UI) CreateHouseForm(w http.ResponseWriter, r *http.Request) {
	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(CreateHouse(), opts...).ServeHTTP(w, r)
}

func (u *UI) HandleCreateHouse(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to decode form values: %v\n", err), 500)
		return
	}

	address := r.Form.Get("address")

	_, err = u.cfg.DB.CreateHouse(r.Context(), database.Text(address))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to create a `House`. Err:%v\n", err), 500)
		return
	}

	http.Redirect(w, r, "/houses", http.StatusSeeOther)
}

func (u *UI) HomePage(w http.ResponseWriter, r *http.Request) {
	dbHouses, err := u.cfg.DB.GetAllHouses(r.Context())
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
		c.Render(r.Context(), w)
		return
	}

	c := HouseView(h)
	c.Render(r.Context(), w)
}

func (u *UI) CreateFlats(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed parsing house_id from the URL: %v\n", err), 500)
		return
	}
	fmt.Println("houseID: ", houseID)
	err = r.ParseForm()
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to parse form: %v\n", err), 500)
		return
	}
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed parsing start of range: %v\n", err), 500)
		return
	}
	to, err := strconv.Atoi(r.Form.Get("to"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed parsing end of range: %v\n", err), 500)
		return
	}

	args := make([]database.CreateFlatsParams, to)
	for i := from - 1; i < to; i++ {
		args[i] = database.CreateFlatsParams{
			HouseID:    houseID,
			FlatNumber: int32(i + 1),
		}
	}

	// fmt.Println(args)

	_, err = u.cfg.DB.CreateFlats(r.Context(), args)
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed to create `[]Flat`: %v\n", err), 500)
		return
	}

	http.Redirect(w, r, "/homes/", http.StatusSeeOther)
}
