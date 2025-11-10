package models

import "github.com/google/uuid"

type FlatOwner struct {
	FlatID uuid.UUID
	Data   Resident
}

type Resident struct {
	ResidentID  uuid.UUID
	FlatID      uuid.UUID
	PhoneNumber string
	Fio         string
}
