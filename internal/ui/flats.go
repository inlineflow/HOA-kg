package ui

import (
	"errors"
	"fmt"
	"hypermedia/internal/database"
	"hypermedia/internal/ui/pages"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

func (u *UI) FlatsCreate(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, fmt.Errorf("Failed parsing house_id from the URL: %v\n", err), 500)
		return
	}

	var opts []func(*templ.ComponentHandler)
	if r.Header.Get("HX-Request") != "" {
		opts = append(opts, templ.WithFragments("partial"))
	}

	templ.Handler(pages.CreateFlats(houseID), opts...).ServeHTTP(w, r)
}

func (u *UI) HandleCreateFlat(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	err = r.ParseForm()
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	err = u.parseFormAndValidateFlatNumber(r, houseID)
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	flatNumber, err := strconv.Atoi(r.Form.Get("flat_number"))
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	flatOwnerName := r.Form.Get("flat_owner_name")
	if flatOwnerName == "" {
		HandleError(w, r, errors.New("Must provide flat owner name"), 400)
		return
	}

	flatOwnerPhoneNumber := r.Form.Get("flat_owner_phone_number")
	if flatOwnerPhoneNumber == "" {
		HandleError(w, r, errors.New("Must provide flat owner phone number"), 400)
		return
	}

	flatArea, err := strconv.Atoi(r.Form.Get("flat_area"))
	if err != nil {
		HandleError(w, r, errors.New("Must provide flat area"), 400)
		return
	}

	err = u.cfg.DB.CreateFlat(r.Context(), u.cfg.Pool, database.CreateFlatParams{
		FlatParams: database.InsertFlatParams{
			FlatNumber: int32(flatNumber),
			HouseID:    houseID,
			Area:       int32(flatArea),
		},
		FlatOwnerResidentParams: database.ResidentInfo{
			Fio:         flatOwnerName,
			PhoneNumber: flatOwnerPhoneNumber,
		},
	})

	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/houses/%s/flats", houseID), http.StatusSeeOther)
}
