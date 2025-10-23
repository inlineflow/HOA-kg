package models

import (
	"hypermedia/internal/database"
)

func ToHouseVM(dbHouse database.House) House {
	return House{
		ID:      dbHouse.HouseID,
		Address: database.PgStringToString(dbHouse.Address),
	}
}

func ToFlatVM(dbFlat database.Flat) Flat {
	return Flat{
		ID:         dbFlat.FlatID,
		HouseID:    dbFlat.HouseID,
		FlatNumber: int(dbFlat.FlatNumber),
	}
}

func Map[T any, V any](source []T, converter func(T) V) []V {
	if source == nil {
		return []V{}
	}

	result := make([]V, len(source))
	for i, v := range source {
		result[i] = converter(v)
	}

	return result
}
