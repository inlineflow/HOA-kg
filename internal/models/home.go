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
