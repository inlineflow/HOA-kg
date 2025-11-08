package ui

import (
	"context"
	"fmt"
	"hypermedia/internal/database"
	"hypermedia/internal/ui/pages"
	"net/http"
	"strconv"

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
	err = u.parseFormAndValidateFlatNumber(r, houseID)
	c := pages.CreateFlatsFormField(
		pages.CreateFlatsFormFieldProps{
			LabelMsg:      m["flat_number"],
			ID:            "flat_number",
			Error:         err,
			InputType:     "number",
			HouseID:       houseID,
			Value:         r.Form.Get("flat_number"),
			ToBeValidated: true,
		},
	)

	c.Render(context.Background(), w)

}

func (u *UI) parseFormAndValidateFlatNumber(r *http.Request, houseID uuid.UUID) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	flatNumber, err := strconv.Atoi(r.Form.Get("flat_number"))
	if err != nil {
		return err
	}

	exists, err := u.cfg.DB.FlatNumberExistsInHouse(r.Context(), database.FlatNumberExistsInHouseParams{
		FlatNumber: int32(flatNumber),
		HouseID:    houseID,
	})

	if exists {
		return fmt.Errorf("Квартира с номером: %d уже существует", flatNumber)
	}

	return nil
}
