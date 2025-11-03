-- name: CreatePaymentPlan :one
INSERT INTO payment_plan(house_id, period_id, date_from, date_to, monthly_amount)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPaymentPlansByHouseID :many
SELECT * FROM payment_plan
where house_id = $1;

-- name: GetLatestPaymentPlanByHouseID :one
SELECT * FROM payment_plan
WHERE house_id = $1
AND date_to IS NULL;

-- name: GetPeriods :many
select * from period;
