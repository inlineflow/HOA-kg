package models

import (
	"hypermedia/internal/database"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type Period struct {
	PeriodID    uuid.UUID
	Tag         string
	Description string
	Group       string
}

func ToPeriodVM(dbPeriod database.Period) Period {
	return Period{
		PeriodID:    dbPeriod.PeriodID,
		Tag:         dbPeriod.Tag,
		Description: dbPeriod.Description,
		Group:       dbPeriod.PeriodGroup,
	}
}

type PaymentPlan struct {
	PaymentPlanID uuid.UUID
	HouseID       uuid.UUID
	PeriodID      uuid.UUID
	DateFrom      *time.Time
	DateTo        *time.Time
	MonthlyAmount decimal.Decimal
}

func ToPaymentPlanVM(dbPP database.PaymentPlan) PaymentPlan {
	return PaymentPlan{
		PaymentPlanID: dbPP.PaymentPlanID,
		HouseID:       dbPP.HouseID,
		PeriodID:      dbPP.PeriodID,
		DateFrom:      database.PgTimestampToTime(dbPP.DateFrom),
		DateTo:        database.PgTimestampToTime(dbPP.DateTo),
		MonthlyAmount: dbPP.MonthlyAmount,
	}
}
