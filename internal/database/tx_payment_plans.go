package database

import (
	"context"
	"fmt"
)

func (q *Queries) CreatePaymentPlan(ctx context.Context, arg InsertPaymentPlanParams) (PaymentPlan, error) {
	err := q.EndLastPlan(ctx, arg.HouseID)
	if err != nil {
		return PaymentPlan{}, fmt.Errorf("Failed to end last plan: %v\n", err)
	}

	pp, err := q.InsertPaymentPlan(ctx, arg)
	if err != nil {
		return PaymentPlan{}, fmt.Errorf("Failed to create a payment plan: %v\n", err)
	}

	return pp, nil
}
