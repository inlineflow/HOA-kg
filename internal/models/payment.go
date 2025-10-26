package models

import (
	"hypermedia/internal/database"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type PaymentPlan struct {
	PaymentPlanID uuid.UUID
	HouseID       uuid.UUID
	PeriodID      uuid.UUID
	DateFrom      time.Time
	DateTo        time.Time
	MonthlyAmount decimal.Decimal
}

func ToPaymentPlanVM(dbPP database.PaymentPlan) {

}
