package models

import "github.com/google/uuid"

type House struct {
	ID      uuid.UUID
	Address string
}

type Flat struct {
	ID         uuid.UUID
	FlatNumber int
}
