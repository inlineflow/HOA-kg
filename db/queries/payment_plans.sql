-- name: InsertPaymentPlan :one
INSERT INTO payment_plan(house_id, period_id, date_from, date_to, monthly_amount)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: EndLastPlan :exec
UPDATE payment_plan
SET date_to = NOW()
WHERE house_id = $1
AND date_to IS NULL;

-- name: GetPaymentPlansByHouseID :many
SELECT sqlc.embed(pp), sqlc.embed(p) FROM payment_plan pp
join period p on pp.period_id = p.period_id
where house_id = $1
order by pp.date_to DESC;

-- name: GetLatestPaymentPlanByHouseID :one
SELECT sqlc.embed(pp), sqlc.embed(p) FROM payment_plan pp
join period p on pp.period_id = p.period_id
WHERE house_id = $1
AND date_to IS NULL;

-- name: GetPeriods :many
select * from period;

