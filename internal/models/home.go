package models

import "github.com/google/uuid"

type House struct {
	ID      uuid.UUID
	Address string
}

type Appartment struct {
	ID         uuid.UUID
	FlatNumber int
}
