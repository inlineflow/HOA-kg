package ui

import (
	"context"
	"fmt"
	"hypermedia/internal/database"
	"hypermedia/internal/ui/pages"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

func getValidate() map[string]string {
	return map[string]string{
		"flat_number":             "Номер квартиры",
		"flat_owner_name":         "Имя владельца",
		"flat_owner_phone_number": "Номер телефона владельца",
		"flat_owner_area":         "Площадь",
	}
}

func (u *UI) ValidateFlatNumber(w http.ResponseWriter, r *http.Request) {
	houseID, err := uuid.Parse(r.PathValue("house_id"))
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	m := getValidate()
	err = r.ParseForm()
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	flatNumber, err := strconv.Atoi(r.Form.Get("flat_number"))
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	exists, err := u.cfg.DB.FlatNumberExistsInHouse(r.Context(), database.FlatNumberExistsInHouseParams{
		FlatNumber: int32(flatNumber),
		HouseID:    houseID,
	})

	if err != nil {
		c := pages.CreateFlatsFormField(
			pages.CreateFlatsFormFieldProps{
				LabelMsg:  m["flat_number"],
				ID:        "flat_number",
				Error:     err,
				InputType: "number",
				HouseID:   houseID,
				Value:     r.Form.Get("flat_number"),
			},
		)
		c.Render(context.Background(), w)
		return
	}

	err = nil
	if exists {
		err = fmt.Errorf("Квартира с номером: %d уже существует", flatNumber)
	}

	c := pages.CreateFlatsFormField(
		pages.CreateFlatsFormFieldProps{
			LabelMsg:  m["flat_number"],
			ID:        "flat_number",
			Error:     err,
			InputType: "number",
			HouseID:   houseID,
			Value:     r.Form.Get("flat_number"),
		},
	)

	c.Render(context.Background(), w)

}
