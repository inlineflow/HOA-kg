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
	Period        Period
	DateFrom      *time.Time
	DateTo        *time.Time
	MonthlyAmount decimal.Decimal
}

func ToPaymentPlanVM(dbPP database.GetPaymentPlansByHouseIDRow) PaymentPlan {
	return PaymentPlan{
		PaymentPlanID: dbPP.PaymentPlan.PaymentPlanID,
		HouseID:       dbPP.PaymentPlan.HouseID,
		Period:        ToPeriodVM(dbPP.Period),
		DateFrom:      database.PgTimestampToTime(dbPP.PaymentPlan.DateFrom),
		DateTo:        database.PgTimestampToTime(dbPP.PaymentPlan.DateTo),
		MonthlyAmount: dbPP.PaymentPlan.MonthlyAmount,
	}
}
