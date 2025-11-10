package models

import (
	"hypermedia/internal/database"

	"github.com/google/uuid"
)

type Flat struct {
	ID            uuid.UUID
	HouseID       uuid.UUID
	FlatNumber    int
	Area          int
	FlatOwnerName string
}

func ToFlatVM(arg database.GetFlatsForHouseRow) Flat {
	return Flat{
		ID:            arg.Flat.FlatID,
		HouseID:       arg.Flat.HouseID,
		FlatNumber:    int(arg.Flat.FlatNumber),
		Area:          int(arg.Flat.Area),
		FlatOwnerName: arg.Resident.Fio,
	}
}
