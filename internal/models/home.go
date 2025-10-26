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
	ID         uuid.UUID
	HouseID    uuid.UUID
	FlatNumber int
}

func ToFlatVM(dbFlat database.Flat) Flat {
	return Flat{
		ID:         dbFlat.FlatID,
		HouseID:    dbFlat.HouseID,
		FlatNumber: int(dbFlat.FlatNumber),
	}
}
