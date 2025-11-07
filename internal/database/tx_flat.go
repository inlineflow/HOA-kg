package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResidentInfo struct {
	PhoneNumber string
	Fio         string
}

type CreateFlatParams struct {
	FlatParams              InsertFlatParams
	FlatOwnerResidentParams ResidentInfo
}

func (q *Queries) CreateFlat(ctx context.Context, pool *pgxpool.Pool, arg CreateFlatParams) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)
	flat, err := qtx.InsertFlat(ctx, arg.FlatParams)
	if err != nil {
		return err
	}

	flatOwnerResident, err := qtx.InsertResident(ctx, InsertResidentParams{
		FlatID:      flat.FlatID,
		PhoneNumber: arg.FlatOwnerResidentParams.PhoneNumber,
		Fio:         arg.FlatOwnerResidentParams.Fio,
	})
	if err != nil {
		return err
	}

	_, err = qtx.InsertFlatOwner(ctx, InsertFlatOwnerParams{
		FlatID:     flat.FlatID,
		ResidentID: flatOwnerResident.ResidentID,
	})

	tx.Commit(ctx)
	return nil

}
