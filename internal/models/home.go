package models

import (
	"hypermedia/internal/database"

	"github.com/google/uuid"
)

type House struct {
	ID      uuid.UUID
	Address string
}

func ToHouseVM(dbHouse database.House) House {
	return House{
		ID:      dbHouse.HouseID,
		Address: dbHouse.Address,
	}
}

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
