package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) CreatePaymentPlan(ctx context.Context, pool *pgxpool.Pool, arg InsertPaymentPlanParams) (PaymentPlan, error) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return PaymentPlan{}, fmt.Errorf("Failed to start transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)
	qtx := q.WithTx(tx)

	err = qtx.EndLastPlan(ctx, arg.HouseID)
	if err != nil {
		return PaymentPlan{}, fmt.Errorf("Failed to end last plan: %v\n", err)
	}

	pp, err := qtx.InsertPaymentPlan(ctx, arg)
	if err != nil {
		return PaymentPlan{}, fmt.Errorf("Failed to create a payment plan: %v\n", err)
	}

	tx.Commit(ctx)
	return pp, nil
}
